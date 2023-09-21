package service

import (
	"context"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/usecase"
	fileusecase "harmoni/internal/usecase/file"

	"go.uber.org/zap"
)

type FileService struct {
	fc     *fileusecase.FileUseCase
	uc     *usecase.UserUseCase
	logger *zap.SugaredLogger
}

func NewFileService(
	fc *fileusecase.FileUseCase,
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

func (s *FileService) GetFileContent(ctx context.Context, req *fileentity.GetFileContentRequest) (*fileentity.GetFileContentResponse, error) {
	content, err := s.fc.GetFileContent(ctx, req.FilePath)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &fileentity.GetFileContentResponse{Content: content}, nil
}

func (s *FileService) UploadObject(ctx context.Context, req *fileentity.UploadObjectRequest) (*fileentity.UploadObjectResponse, error) {
	file, err := s.fc.Save(ctx, &fileentity.File{
		Name: req.FileName,
		Size: uint64(req.Size),
	}, req.Content)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	link, err := s.fc.GetFileLink(ctx, file.FileID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &fileentity.UploadObjectResponse{Location: link}, nil
}

func (s *FileService) IsObjectUploaded(ctx context.Context, req *fileentity.IsObjectUploadedRequest) (*fileentity.IsObjectUploadedResponse, error) {
	location, err := s.fc.IsObjectUploaded(ctx, req.Hash)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.IsObjectUploadedResponse{Location: location}, nil
}

func (s *FileService) UploadPrepare(ctx context.Context, req *fileentity.UploadPrepareRequest) (*fileentity.CreateUploadSessionResponse, error) {
	credential, err := s.fc.UploadPrepare(ctx, req.Key, req.MD5, req.UserID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.CreateUploadSessionResponse{Credential: credential}, nil
}

func (s *FileService) UploadPart(ctx context.Context, req *fileentity.UploadMultipartRequest) (*fileentity.UploadMultipartResponse, error) {
	objectID, err := s.fc.UploadPart(ctx, req.UploadID, req.Key, req.PartNumber, req.UserID, req.Content, uint64(req.Size))
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.UploadMultipartResponse{Etag: objectID}, nil
}

func (s *FileService) UploadComplete(ctx context.Context, req *fileentity.UploadMultipartCompleteRequest) (*fileentity.UploadMultipartCompleteResponse, error) {
	eTag, location, err := s.fc.UploadComplete(ctx, req.UploadID, req.Key, req.UserID, req.FileParts)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.UploadMultipartCompleteResponse{Etag: eTag, Location: location}, nil
}

func (s *FileService) ListParts(ctx context.Context, req *fileentity.ListPartsRequest) (*fileentity.ListPartsResponse, error) {
	parts, err := s.fc.ListParts(ctx, req.UploadID,
		req.Key, req.UserID, req.MaxParts, req.PartNumberMarker)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.ListPartsResponse{Parts: parts}, nil
}

func (s *FileService) AbortMultipartUpload(ctx context.Context, req *fileentity.AbortMultipartUploadRequest) (*fileentity.AbortMultipartUploadResponse, error) {
	err := s.fc.AbortMultipartUpload(ctx, req.UploadID, req.Key, req.UserID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}
	return &fileentity.AbortMultipartUploadResponse{}, nil
}
