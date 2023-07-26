package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"harmoni/internal/conf"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type FileUseCase struct {
	appConf        *conf.App
	fileConf       *conf.FileStorage
	fileRepository fileentity.FileRepository
	logger         *zap.SugaredLogger
}

// store file in local for now
func NewFileUseCase(
	appConf *conf.App,
	fileConf *conf.FileStorage,
	fileRepository fileentity.FileRepository,
	logger *zap.SugaredLogger,
) *FileUseCase {
	if fileConf.Type == string(fileentity.LocalStorageType) {
		createFolderIfNotExists(fileConf.Local.Path)
	}

	return &FileUseCase{
		appConf:        appConf,
		fileConf:       fileConf,
		fileRepository: fileRepository,
		logger:         logger,
	}
}

const (
	basic_avatar = "noface.jpg"
)

func createFolderIfNotExists(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}

func (u *FileUseCase) Save(ctx context.Context, file *fileentity.File, fileContent []byte) (*fileentity.File, error) {
	ext := filepath.Ext(file.Name)
	var path string
	switch u.fileConf.Type {
	case string(fileentity.LocalStorageType):
		path = u.fileConf.Local.Path
	}

	h := md5.Sum(fileContent)
	file.Hash = hex.EncodeToString(h[:])
	file.Path = filepath.Join(path, file.Hash+ext)
	file.Ext = ext
	file.Name = file.Hash + file.Ext

	var err error
	file, err = u.fileRepository.Save(ctx, file)
	if err != nil {
		return nil, err
	}

	switch u.fileConf.Type {
	case string(fileentity.LocalStorageType):
		file, err := os.Create(file.Path)
		if err != nil {
			return nil, errorx.InternalServer(reason.ServerError).WithError(err).WithStack()
		}
		defer file.Close()
		written, err := file.Write(fileContent)
		if written != len(fileContent) {
			u.logger.Errorf("want write %v bytes, but %v bytes", len(fileContent), written)
		}
		if err != nil {
			return nil, errorx.InternalServer(reason.ServerError).WithError(err).WithStack()
		}
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

	var link string
	switch u.fileConf.Type {
	case string(fileentity.LocalStorageType):
		link, err = url.JoinPath(u.appConf.BaseURL, file.Path)
		if err != nil {
			return "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return link, nil
	}

	return link, nil
}

func (u *FileUseCase) GetFileContent(ctx context.Context, filename string) ([]byte, error) {
	file, err := u.fileRepository.GetByFileHash(ctx, filename)
	if err != nil {
		return nil, err
	}

	var content []byte
	switch u.fileConf.Type {
	case string(fileentity.LocalStorageType):
		return u.readFromLocal(ctx, file.Path)
	}

	return content, nil
}

func (u *FileUseCase) readFromLocal(ctx context.Context, path string) ([]byte, error) {
	// store file in local for now
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return content, nil
}
