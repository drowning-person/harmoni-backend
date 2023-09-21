package post

import (
	"context"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/entity/user"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Status    int32   `gorm:"not null"`
	PostID    int64   `gorm:"uniqueIndex"`
	AuthorID  int64   `gorm:"index"`
	Title     string  `gorm:"type:varchar(128)"`
	Content   string  `gorm:"type:text"`
	LikeCount int64   `gorm:"not null"`
	TagIDs    []int64 `gorm:"-"`
}

func (Post) TableName() string {
	return "post"
}

type PostInfo struct {
	PostBasicInfo
	CollectCount int64 `json:"collect_count"`
	Collected    bool  `json:"collected"`
}

type PostBasicInfo struct {
	Liked        bool                `json:"liked"`
	Status       int32               `json:"status"`
	LikeCount    int64               `json:"like_count"`
	CommentCount int64               `json:"comment_count"`
	PostID       int64               `json:"post_id,string"`
	User         *user.UserBasicInfo `json:"user_info"`
	Tags         []tagentity.TagInfo `json:"tags"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (p *Post) ToBasic() PostBasicInfo {
	return PostBasicInfo{
		User: &user.UserBasicInfo{
			UserID: p.AuthorID,
		},
		Status:    p.Status,
		PostID:    p.PostID,
		Title:     p.Title,
		Content:   p.Content,
		LikeCount: p.LikeCount,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *Post) ToInfo() PostInfo {
	pd := PostInfo{
		PostBasicInfo: p.ToBasic(),
	}

	return pd
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Delete(ctx context.Context, postID int64) error
	GetBasicInfoByPostID(ctx context.Context, postID int64) (*Post, bool, error)
	GetByUserID(ctx context.Context, userID int64, queryCond *PostQuery) (paginator.Page[Post], error)
	GetByUserIDs(ctx context.Context, userID []int64, queryCond *PostQuery) (paginator.Page[Post], error)
	GetByPostID(ctx context.Context, postID int64) (*Post, bool, error)
	GetPostsByTagID(ctx context.Context, tagID int64) ([]Post, error)
	BatchByIDs(ctx context.Context, postIDs []int64) ([]Post, error)
	GetLikeCount(ctx context.Context, postID int64) (int64, bool, error)
	UpdateLikeCount(ctx context.Context, postID int64, count int64) error
	GetPage(ctx context.Context, queryCond *PostQuery) (paginator.Page[Post], error)
}
