package file

import (
	"context"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/entity/unique"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ fileentity.FileRepository = (*FileRepo)(nil)

type FileRepo struct {
	db           *gorm.DB
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewFileRepository(
	db *gorm.DB,
	uniqueIDRepo unique.UniqueIDRepo,
	logger *zap.SugaredLogger,
) *FileRepo {
	return &FileRepo{
		db:           db,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/file"),
	}
}

func (r *FileRepo) Save(ctx context.Context, file *fileentity.File) (*fileentity.File, error) {
	var err error
	file.FileID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return nil, err
	}

	err = r.db.Create(file).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return file, nil
}

func (r *FileRepo) GetByFileID(ctx context.Context, fileID int64) (*fileentity.File, error) {
	var file fileentity.File
	if err := r.db.Where("file_id = ?", fileID).First(&file).Error; err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return &file, nil
}

func (r *FileRepo) GetByFileHash(ctx context.Context, fileHash string) (*fileentity.File, error) {
	var file fileentity.File
	if err := r.db.Where("hash = ?", fileHash).First(&file).Error; err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return &file, nil
}
