package comment

import (
	"context"
	"harmoni/internal/entity/paginator"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID int64 `gorm:"uniqueIndex"`
	ObjectID  int64
	AuthorID  int64
	ParentID  int64
	RootID    int64
	Content   string `gorm:"type:varchar(2048)"`
	LikeCount int64  `gorm:"not null"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentDetail struct {
	CommentID int64     `json:"cid,string,omitempty"`
	ObjectID  int64     `json:"oid,string,omitempty"`
	AuthorID  int64     `json:"aid,string,omitempty"`
	ParentID  int64     `json:"pid,string,omitempty"`
	RootID    int64     `json:"rid,string,omitempty"`
	Content   string    `json:"content,omitempty"`
	LikeCount int64     `json:"like_count,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func CommentToCommentDetail(c *Comment) CommentDetail {
	return CommentDetail{
		CommentID: c.CommentID,
		ObjectID:  c.ObjectID,
		AuthorID:  c.AuthorID,
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		LikeCount: c.LikeCount,
	}
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByCommentID(ctx context.Context, commentID int64) (*Comment, bool, error)
	GetLikeCount(ctx context.Context, commentID int64) (int64, bool, error)
	UpdateLikeCount(ctx context.Context, commentID int64, count int64) error
	GetPage(ctx context.Context, commentQuery *CommentQuery) (paginator.Page[Comment], error)
}
