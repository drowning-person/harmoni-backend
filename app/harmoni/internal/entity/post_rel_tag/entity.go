package postreltag

import (
	"context"
	"harmoni/app/harmoni/internal/entity"
	"harmoni/app/harmoni/internal/entity/post"
	"harmoni/app/harmoni/internal/entity/tag"
)

const TableName = "post_tags"

type PostTag struct {
	entity.BaseModel
	PostTagID int64 `gorm:"uniqueIndex"`
	PostID    int64 `gorm:"index"`
	TagID     int64 `gorm:"index"`
}

type PostRelTagRepository interface {
	AssociateTags(ctx context.Context, postID int64, tagIDs []int64) error
	GetTagsByPostID(ctx context.Context, postID int64) ([]tag.Tag, error)
	GetPostsByTagID(ctx context.Context, tagID int64) ([]post.Post, error)
	RemoveTagsFromPost(ctx context.Context, postID int64, tagIDs []int64) error
	RemoveAllTagsFromPost(ctx context.Context, postID int64) error
}

func (PostTag) TableName() string {
	return TableName
}
