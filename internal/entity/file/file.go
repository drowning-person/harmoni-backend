package file

import (
	"context"
	"harmoni/internal/entity"
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
	Size   uint64 `gorm:"not null"`
}

// 自定义表名
func (File) TableName() string {
	return TableName
}

type FileRepository interface {
	Save(ctx context.Context, file *File) (*File, error)
	GetByFileID(ctx context.Context, fileID int64) (*File, error)
	GetByFileHash(ctx context.Context, fileHash string) (*File, error)
}
