package service

import (
	"context"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type FileService struct {
	fc     *usecase.FileUseCase
	uc     *usecase.UserUseCase
	logger *zap.SugaredLogger
}

func NewFileService(
	fc *usecase.FileUseCase,
	uc *usecase.UserUseCase,
	logger *zap.SugaredLogger,
) *FileService {
	return &FileService{
		fc:     fc,
		uc:     uc,
		logger: logger.With("module", "service/file"),
	}
}

func (s *FileService) UploadAvatar(ctx context.Context, req *fileentity.AvatarUploadRequest) (*fileentity.AvatarUploadResponse, error) {
	file, err := s.fc.Save(ctx, &fileentity.File{
		Name: req.FileName,
		Size: uint64(req.Size),
	}, req.Content)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	err = s.uc.SetAvatar(ctx, req.UserID, file.FileID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	link, err := s.fc.GetFileLink(ctx, file.FileID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &fileentity.AvatarUploadResponse{Link: link}, nil
}
