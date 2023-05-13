package tag

import (
	"context"
	"harmoni/internal/entity/paginator"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	TagID        int64  `gorm:"not null;uniqueIndex"`
	TagName      string `gorm:"not null;type:varchar(128);unique"`
	Introduction string `gorm:"not null;type:varchar(256)"`
	FollowCount  int64  `gorm:"not null;default:0"`
}

type TagInfo struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}

type TagDetail struct {
	TagID        int64     `json:"tag_id,omitempty"`
	TagName      string    `json:"tag_name,omitempty"`
	Introduction string    `json:"introduction,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	FollowCount  int64     `json:"follow_count"`
}

func (*Tag) TableName() string {
	return "tag"
}

func ConvertTagToDisplay(t *Tag) TagInfo {
	return TagInfo{
		TagID:   t.TagID,
		TagName: t.TagName,
	}
}

func ConvertTagToDetailDisplay(t *Tag) TagDetail {
	return TagDetail{
		TagID:        t.TagID,
		TagName:      t.TagName,
		Introduction: t.Introduction,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		FollowCount:  t.FollowCount,
	}
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	GetByTagID(ctx context.Context, tagID int64) (*Tag, bool, error)
	GetByTagIDs(ctx context.Context, tagID []int64) ([]Tag, error)
	GetByTagName(ctx context.Context, tagName string) (*Tag, bool, error)
	GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[Tag], error)
}
