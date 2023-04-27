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
}

type CommentDetail struct {
	CommentID int64     `json:"cid,string"`
	ObjectID  int64     `json:"oid,string"`
	AuthorID  int64     `json:"aid,string"`
	ParentID  int64     `json:"pid,string"`
	RootID    int64     `json:"rid,string"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	}
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByCommentID(ctx context.Context, commentID int64) (*Comment, bool, error)
	GetPage(ctx context.Context, commentQuery *CommentQuery) (paginator.Page[Comment], error)
}
