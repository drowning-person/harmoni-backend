package file

import (
	"context"
	"harmoni/internal/entity"
	"harmoni/internal/pkg/filesystem/upload"
	"time"
)

const (
	TableName = "file"
)

type StorageType string

const (
	LocalStorageType StorageType = "local"
)

type File struct {
	entity.BaseModelWithNoSoftDelete
	FileID int64  `gorm:"not null;uniqueIndex"`
	Name   string `gorm:"not null"`
	Path   string `gorm:"not null"`
	Ext    string `gorm:"not null"`
	Hash   string `gorm:"not null"`
	ETag   string `gorm:"not null"`
	Size   uint64 `gorm:"not null"`
}

// 自定义表名
func (File) TableName() string {
	return TableName
}

type FileList []*File

func (f FileList) ToMap() map[int64]*File {
	m := make(map[int64]*File)
	for i, file := range f {
		m[file.FileID] = f[i]
	}
	return m
}

type FileRepository interface {
	Save(ctx context.Context, file *File) (*File, error)
	GetByFileID(ctx context.Context, fileID int64) (*File, error)
	ListByFileID(ctx context.Context, fileIDs []int64) (FileList, error)
	GetByFileHash(ctx context.Context, fileHash string) (*File, error)
	UpdatePartInfo(ctx context.Context, file *File) error

	CreateUploadSession(ctx context.Context, uploadSession *upload.UploadSession, ttl time.Duration) error
	GetUploadSession(ctx context.Context, uploadID string) (*upload.UploadSession, error)
	DeleteUploadSession(ctx context.Context, uploadID string) error
}
