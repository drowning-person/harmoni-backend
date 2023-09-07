package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"harmoni/internal/conf"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/filesystem"
	"harmoni/internal/pkg/filesystem/driver"
	"harmoni/internal/pkg/filesystem/fsctx"
	"harmoni/internal/pkg/filesystem/response"
	"harmoni/internal/pkg/filesystem/upload"
	"harmoni/internal/pkg/reason"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type FileUseCase struct {
	appConf        *conf.App
	rdb            redis.UniversalClient
	fs             *filesystem.FileSystem
	fileConf       *conf.FileStorage
	fileRepository fileentity.FileRepository
	logger         *zap.SugaredLogger
}

// store file in local for now
func NewFileUseCase(
	appConf *conf.App,
	rdb redis.UniversalClient,
	fileConf *conf.FileStorage,
	filesystem *filesystem.FileSystem,
	fileRepository fileentity.FileRepository,
	logger *zap.SugaredLogger,
) *FileUseCase {

	return &FileUseCase{
		appConf:        appConf,
		rdb:            rdb,
		fileConf:       fileConf,
		fs:             filesystem,
		fileRepository: fileRepository,
		logger:         logger,
	}
}

const (
	basic_avatar = "noface.jpg"
)

func (u *FileUseCase) Save(ctx context.Context, file *fileentity.File, fileContent io.ReadSeekCloser) (*fileentity.File, error) {
	defer u.fs.CleanHooks("")
	ext := filepath.Ext(file.Name)
	data := make([]byte, file.Size)
	n, err := fileContent.Read(data)

	if err != nil {
		return nil, err
	} else if n != int(file.Size) {
		return nil, errorx.BadRequest(reason.UploadFileSizeIncorrect).WithError(fmt.Errorf("want %d bytes but %d bytes", file.Size, n))
	}

	h := md5.Sum(data)
	hash := hex.EncodeToString(h[:])
	if file, err := u.fileRepository.GetByFileHash(ctx, hash); err != nil {
		return nil, err
	} else if file != nil {
		return nil, errorx.BadRequest(reason.UploadFileSliceUploaded)
	}
	fileContent.Seek(0, io.SeekStart)
	u.fs.Use("BeforeUpload", u.HookValidateFile)
	fileStream := &fsctx.FileStream{
		Mode: fsctx.Overwrite,
		Hash: hash,
		Ext:  ext,
		Name: hash + ext,
		File: fileContent,
		Size: uint64(len(data)),
	}

	eTag, err := u.fs.Upload(ctx, fileStream)
	if err != nil {
		return nil, err
	}
	fileStream.ETag = eTag
	file, err = u.AddFile(ctx, fileStream)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (u *FileUseCase) GetFileLink(ctx context.Context, fileID int64) (string, error) {
	if fileID == 0 {
		return filepath.Join(u.appConf.BaseURL, basic_avatar), nil
	}
	file, err := u.fileRepository.GetByFileID(ctx, fileID)
	if err != nil {
		return "", err
	}
	rootURL := ""
	if u.fs.Policy.Type == string(fileentity.LocalStorageType) {
		rootURL = u.appConf.BaseURL
	}
	link, err := u.fs.Handler.Source(ctx, rootURL, file.Path, 0, false, 0)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (u *FileUseCase) GetFileContent(ctx context.Context, filepath string) (response.RSCloser, error) {
	content, err := u.fs.Handler.Get(ctx, filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errorx.NotFound(reason.FileNotFound)
		}
		return nil, err
	}
	return content, nil
}

func (u *FileUseCase) IsObjectUploaded(ctx context.Context, hash string) (string, error) {
	file, err := u.fileRepository.GetByFileHash(ctx, hash)
	if err != nil {
		return "", err
	} else if file == nil {
		return "", nil
	}

	return u.GetFileLink(ctx, file.FileID)
}

func (u *FileUseCase) UploadPrepare(ctx context.Context, key string, md5 string, userID int64) (*upload.UploadCredential, error) {
	defer u.fs.CleanHooks("")
	if file, err := u.fileRepository.GetByFileHash(ctx, md5); err != nil {
		return nil, err
	} else if file != nil {
		return nil, errorx.BadRequest(reason.UploadFileSliceUploaded)
	}

	// 获取相关有效期设置
	callBackSessionTTL := u.fileConf.UploadSessionTimeout
	file := &fsctx.FileStream{
		IsPart: true,
		Name:   key,
		Mode:   fsctx.Nop,
		Ext:    filepath.Ext(key),
		File:   io.NopCloser(strings.NewReader("")),
	}
	callbackKey := uuid.Must(uuid.NewRandom()).String()

	// 创建占位的文件，同时校验文件信息
	if callbackKey != "" {
		file.UploadID = callbackKey
	}

	u.fs.Use("BeforeUpload", u.HookValidateFile)

	// 验证文件规格
	if _, err := u.fs.Upload(ctx, file); err != nil {
		return nil, err
	}

	uploadSession := &upload.UploadSession{
		UploadID: callbackKey,
		UID:      uint(userID),
		Policy:   u.fs.Policy,
		Key:      key,
		SavePath: file.SavePath,
	}

	// 获取上传凭证
	credential, err := u.fs.Handler.Token(ctx, int64(callBackSessionTTL), uploadSession, file)
	if err != nil {
		return nil, err
	}

	u.fs.Use("AfterUpload", u.GenericAfterUpload)
	ctx = context.WithValue(ctx, fsctx.IgnoreDirectoryConflictCtx, true)
	if _, err := u.fs.Upload(ctx, file); err != nil {
		return nil, err
	}
	uploadSession.FileID = file.Model.(*fileentity.File).FileID
	// 创建回调会话
	err = u.fileRepository.CreateUploadSession(ctx, uploadSession, callBackSessionTTL)
	if err != nil {
		return nil, err
	}
	// 补全上传凭证其他信息
	credential.Expires = time.Now().Add(callBackSessionTTL).Unix()

	return credential, nil
}

func (u *FileUseCase) validateUploadSession(ctx context.Context, uploadID string, key string, userID int64) (*upload.UploadSession, error) {
	uploadSession, err := u.fileRepository.GetUploadSession(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	if uploadSession.UID != uint(userID) {
		return nil, errorx.BadRequest(reason.UploadFileSizeIncorrect)
	}
	if uploadSession.Key != key {
		return nil, errors.New("upload session expired")
	}

	return uploadSession, nil
}

// UploadPart 处理文件分片上传
func (u *FileUseCase) UploadPart(ctx context.Context, uploadID string, key string, index int, userID int64, data io.ReadCloser, size uint64) (string, error) {
	defer u.fs.CleanHooks("")

	uploadSession, err := u.validateUploadSession(ctx, uploadID, key, userID)
	if err != nil {
		return "", err
	}

	// 查找上传会话创建的占位文件
	file, err := u.fileRepository.GetByFileID(ctx, uploadSession.FileID)
	if err != nil {
		return "", errors.New("upload session expired")
	}

	return u.processChunkUpload(ctx, data, uploadSession, index, file, fsctx.Append, size)
}

func (u *FileUseCase) processChunkUpload(ctx context.Context, data io.ReadCloser, session *upload.UploadSession, index int, file *fileentity.File, mode fsctx.WriteMode, size uint64) (string, error) {
	// 非首个分片时需要允许覆盖
	if index > 0 {
		mode |= fsctx.Overwrite
	}

	fileData := fsctx.FileStream{
		IsPart:     true,
		PartNumber: index,
		File:       data,
		Size:       size,
		Mode:       mode,
		UploadID:   session.UploadID,
		Ext:        filepath.Ext(session.SavePath),
	}

	if file != nil {
		u.fs.Use("AfterUpload", u.HookChunkUploaded(session.UploadID, index))
	}

	// 执行上传
	etag, err := u.fs.Upload(ctx, &fileData)
	if err != nil {
		return "", err
	}

	return etag, nil
}

// UploadComplete merge file parts
// return ETag, link and error
func (u *FileUseCase) UploadComplete(ctx context.Context, uploadID string, key string, userID int64, parts []fileentity.FilePart) (string, string, error) {
	defer u.fs.CleanHooks("")

	uploadSession, err := u.validateUploadSession(ctx, uploadID, key, userID)
	if err != nil {
		return "", "", err
	}

	fileData := fsctx.FileStream{
		IsPart:   true,
		UploadID: uploadSession.UploadID,
		Policy:   uploadSession.Policy,
		Ext:      filepath.Ext(uploadSession.SavePath),
		SavePath: uploadSession.SavePath,
	}

	hash, err := u.fs.Handler.Merge(ctx, &fileData, fileentity.Parts(parts).ToDriver())
	if err != nil {
		if err == driver.ErrInvalidParts {
			return "", "", errorx.BadRequest(reason.UploadInvalidPart).WithMsg(err.Error())
		} else if errors.Is(err, driver.ErrEntityTooSmall) {
			{
				err := u.fs.Handler.DeleteParts(ctx, &fileData)
				if err != nil {
					return "", "", err
				}

				err = u.fileRepository.DeleteUploadSession(ctx, uploadID)
				if err != nil {
					return "", "", err
				}
			}
			return "", "", errorx.BadRequest(reason.UploadEntityTooSmall).WithMsg(err.Error())
		}
		return "", "", err
	}
	err = u.fileRepository.UpdatePartInfo(ctx, &fileentity.File{
		FileID: uploadSession.FileID,
		ETag:   hash,
		Hash:   hash,
		Size:   fileData.Size,
	})
	if err != nil {
		return "", "", err
	}

	rootURL := ""
	if u.fs.Policy.Type == string(fileentity.LocalStorageType) {
		rootURL = u.appConf.BaseURL
	}
	link, err := u.fs.Handler.Source(ctx, rootURL, uploadSession.SavePath, 0, false, 0)
	if err != nil {
		return "", "", err
	}

	err = u.fileRepository.DeleteUploadSession(ctx, uploadID)
	if err != nil {
		return "", "", err
	}

	return hash, link, nil
}

func (u *FileUseCase) ListParts(ctx context.Context, uploadID string, key string, userID int64, maxParts int64, offset int64) ([]driver.Part, error) {
	uploadSession, err := u.validateUploadSession(ctx, uploadID, key, userID)
	if err != nil {
		return nil, err
	}
	fileData := fsctx.FileStream{
		IsPart:   true,
		UploadID: uploadSession.UploadID,
	}

	parts, err := u.fs.Handler.ListParts(ctx, &fileData, maxParts, offset)
	if err != nil {
		if err == driver.ErrInvalidParts {
			return nil, errorx.BadRequest(reason.UploadInvalidPart).WithMsg(err.Error())
		}
		return nil, err
	}

	return parts, nil
}

func (u *FileUseCase) AbortMultipartUpload(ctx context.Context, uploadID string, key string, userID int64) error {
	uploadSession, err := u.validateUploadSession(ctx, uploadID, key, userID)
	if err != nil {
		return err
	}
	fileData := fsctx.FileStream{
		IsPart:   true,
		UploadID: uploadSession.UploadID,
	}

	err = u.fs.Handler.DeleteParts(ctx, &fileData)
	if err != nil {
		if err == driver.ErrInvalidParts {
			return errorx.BadRequest(reason.UploadInvalidPart).WithMsg(err.Error())
		}
		return err
	}

	err = u.fileRepository.DeleteUploadSession(ctx, uploadID)
	if err != nil {
		return err
	}

	return nil
}
