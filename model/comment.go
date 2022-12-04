package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID int64 `gorm:"uniqueIndex"`
	PostID    int64
	AuthorID  int64
	ParentID  int64
	RootID    int64
	Content   string `gorm:"type:varchar(512)"`
}

type CommentDetail struct {
	CommentID int64     `json:"comment_id,string"`
	PostID    int64     `json:"post_id,string"`
	Auther    UserInfo  `json:"author"`
	ParentID  int64     `json:"parent_id,string"`
	RootID    int64     `json:"root_id,string"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CommentToCommentDetail(c *Comment) (*CommentDetail, error) {
	user := UserInfo{UserID: c.AuthorID}
	if err := DB.Table("users").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &CommentDetail{
		CommentID: c.CommentID,
		PostID:    c.PostID,
		Auther:    user,
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

func CommentsToCommentDetails(cs []Comment) ([]CommentDetail, error) {
	cds := make([]CommentDetail, 0, len(cs))
	for _, v := range cs {
		if cd, err := CommentToCommentDetail(&v); err != nil {
			return nil, err
		} else if cd == nil {
			return nil, nil
		} else {
			cds = append(cds, *cd)
		}
	}
	return cds, nil
}
