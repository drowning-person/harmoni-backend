package comment

import (
	"encoding/json"
	"harmoni/app/harmoni/internal/entity/comment"
	"harmoni/app/harmoni/internal/entity/user"
	"harmoni/app/harmoni/internal/pkg/common"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID int64 `gorm:"uniqueIndex"`
	ObjectID  int64
	AuthorID  int64
	ParentID  int64
	RootID    int64
	DialogID  int64
	Content   string `gorm:"type:varchar(2048)"`
	ToUserIDs string `gorm:"type:varchar(2048)"`
	LikeCount int64  `gorm:"not null"`
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) ToDomain() *comment.Comment {
	comm := &comment.Comment{
		CommentID: c.CommentID,
		ObjectID:  c.ObjectID,
		Author: &user.UserBasicInfo{
			UserID: c.AuthorID,
		},
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		LikeCount: c.LikeCount,
	}
	if len(c.ToUserIDs) != 0 {
		comm.ToMembers = make([]*user.UserBasicInfo, 0, len(c.ToUserIDs))
		data := []int64{}
		json.Unmarshal(common.StringToBytes(c.ToUserIDs), &data)
		for _, toUserID := range data {
			comm.ToMembers = append(comm.ToMembers, &user.UserBasicInfo{UserID: toUserID})
		}
	}
	return comm
}

func (c *Comment) FromDomain(ce *comment.Comment) *Comment {
	c.CommentID = ce.CommentID
	c.ObjectID = ce.ObjectID
	c.AuthorID = ce.Author.UserID
	c.ParentID = ce.ParentID
	c.RootID = ce.RootID
	c.Content = ce.Content
	c.CreatedAt = ce.CreatedAt
	c.UpdatedAt = ce.UpdatedAt
	c.LikeCount = ce.LikeCount
	toUserIDs := make([]int64, 0, len(ce.ToMembers))
	for _, toUser := range ce.ToMembers {
		toUserIDs = append(toUserIDs, toUser.UserID)
	}
	data, _ := json.Marshal(&toUserIDs)
	c.ToUserIDs = common.BytesToString(data)
	return c
}
