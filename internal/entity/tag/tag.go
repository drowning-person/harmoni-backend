package tag

import (
	"context"
	"harmoni/internal/entity/paginator"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	TagID        int64
	TagName      string `gorm:"varchar(128);unique"`
	Introduction string `gorm:"varchar(256)"`
}

type TagInfo struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}

type TagDetail struct {
	TagID        int64     `json:"tag_id"`
	TagName      string    `json:"tag_name"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ConvertTagToDisplay(t Tag) TagInfo {
	return TagInfo{
		TagID:   t.TagID,
		TagName: t.TagName,
	}
}

func ConvertTagToDetailDisplay(t Tag) TagDetail {
	return TagDetail{
		TagID:        t.TagID,
		TagName:      t.TagName,
		Introduction: t.Introduction,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	GetByTagID(ctx context.Context, tagID int64) (Tag, error)
	GetByTagName(ctx context.Context, tagName string) (Tag, error)
	GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[Tag], error)
}
