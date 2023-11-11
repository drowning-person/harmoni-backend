package file

import (
	"context"
	"fmt"
	fileentity "harmoni/app/harmoni/internal/entity/file"
	"harmoni/app/harmoni/internal/entity/unique"
	"harmoni/app/harmoni/internal/pkg/filesystem/upload"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ fileentity.FileRepository = (*FileRepo)(nil)

const (
	UploadSessionCachePrefix = "upload_session:"
)

func getUploadSessionKey(uploadID string) string {
	return fmt.Sprintf("%s%s", UploadSessionCachePrefix, uploadID)
}

type FileRepo struct {
	db           *gorm.DB
	rdb          redis.UniversalClient
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewFileRepository(
	db *gorm.DB,
	rdb redis.UniversalClient,
	uniqueIDRepo unique.UniqueIDRepo,
	logger *zap.SugaredLogger,
) *FileRepo {
	return &FileRepo{
		db:           db,
		rdb:          rdb,
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
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NotFound(reason.FileNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return &file, nil
}

func (r *FileRepo) ListByFileID(ctx context.Context, fileIDs []int64) (fileentity.FileList, error) {
	fileList := make([]*fileentity.File, 0, len(fileIDs))
	if err := r.db.Where("file_id IN ?", fileIDs).Find(&fileList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NotFound(reason.FileNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return fileList, nil
}

func (r *FileRepo) GetByFileHash(ctx context.Context, fileHash string) (*fileentity.File, error) {
	var file fileentity.File
	if err := r.db.WithContext(ctx).Where("hash = ?", fileHash).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return &file, nil
}

func (r *FileRepo) UpdatePartInfo(ctx context.Context, file *fileentity.File) error {
	err := r.db.WithContext(ctx).
		Table(file.TableName()).
		Select("hash", "e_tag", "size").
		Where("file_id = ?", file.FileID).
		Updates(file).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *FileRepo) CreateUploadSession(ctx context.Context, uploadSession *upload.UploadSession, exp time.Duration) error {
	cacheKey := getUploadSessionKey(uploadSession.UploadID)
	err := r.rdb.Set(ctx, cacheKey, uploadSession.ToJSON(), exp).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	err = r.rdb.Expire(ctx, cacheKey, exp).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *FileRepo) GetUploadSession(ctx context.Context, uploadID string) (*upload.UploadSession, error) {
	cacheKey := getUploadSessionKey(uploadID)
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errorx.NotFound(reason.UploadSessionNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return (&upload.UploadSession{}).FromJSONString(data), nil
}

func (r *FileRepo) DeleteUploadSession(ctx context.Context, uploadID string) error {
	cacheKey := getUploadSessionKey(uploadID)
	err := r.rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
