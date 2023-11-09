package file

import (
	"context"
	fileentity "harmoni/app/harmoni/internal/entity/file"
	"harmoni/app/harmoni/internal/pkg/filesystem"
	"harmoni/app/harmoni/internal/pkg/filesystem/fsctx"
)

// GenericAfterUpload 文件上传完成后，包含数据库操作
func (u *FileUseCase) GenericAfterUpload(ctx context.Context, fs *filesystem.FileSystem, fileHeader fsctx.FileHeader) error {
	// 向数据库中插入记录
	file, err := u.AddFile(ctx, fileHeader)
	if err != nil {
		return filesystem.ErrInsertFileRecord
	}
	fileHeader.SetModel(file)
	return nil
}

func (u *FileUseCase) AddFile(ctx context.Context, file fsctx.FileHeader) (*fileentity.File, error) {
	uploadInfo := file.Info()
	newFile := &fileentity.File{
		Name: uploadInfo.FileName,
		Size: uploadInfo.Size,
		Hash: uploadInfo.Hash,
		Ext:  uploadInfo.Ext,
		Path: uploadInfo.SavePath,
	}
	var err error
	newFile, err = u.fileRepository.Save(ctx, newFile)
	if err != nil {
		return nil, err
	}

	return newFile, nil
}

// HookValidateFile 一系列对文件检验的集合
func (u *FileUseCase) HookValidateFile(ctx context.Context, fs *filesystem.FileSystem, file fsctx.FileHeader) error {
	fileInfo := file.Info()

	// 验证单文件尺寸
	if !u.ValidateFileSize(ctx, fileInfo.Size) {
		return filesystem.ErrFileSizeTooBig
	}

	// 验证文件名
	if !u.ValidateLegalName(ctx, fileInfo.FileName) {
		return filesystem.ErrIllegalObjectName
	}

	// 验证扩展名
	if !u.ValidateExtension(ctx, fileInfo.FileName) {
		return filesystem.ErrFileExtensionNotAllowed
	}

	return nil

}

// HookChunkUploadFinished 单个分片上传结束后
func (u *FileUseCase) HookChunkUploaded(sessionID string, partNumber int) filesystem.Hook {
	return func(ctx context.Context, fs *filesystem.FileSystem, fileHeader fsctx.FileHeader) error {
		// u.rdb.HSet(ctx, getUploadSessionKey(sessionID), partNumber, "")
		return nil
	}
}
