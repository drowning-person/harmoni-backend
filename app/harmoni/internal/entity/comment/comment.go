package comment

import (
	"context"
	"harmoni/app/harmoni/internal/entity/paginator"
	"harmoni/app/harmoni/internal/entity/user"
	"html"
	"time"
)

type Comment struct {
	CommentID int64                 `json:"cid,string"`
	ObjectID  int64                 `json:"oid,string"`
	Author    *user.UserBasicInfo   `json:"author"`
	ToMembers []*user.UserBasicInfo `json:"toMembers"`
	ParentID  int64                 `json:"pid,string"`
	RootID    int64                 `json:"rid,string"`
	Content   string                `json:"content"`
	LikeCount int64                 `json:"like_count"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Children  []*Comment            `json:"replies"`
}

type CommentList []*Comment

func (l CommentList) ToRooIDMap() map[int64][]*Comment {
	m := map[int64][]*Comment{}
	for i := range l {
		m[l[i].RootID] = append(m[l[i].RootID], l[i])
	}
	return m
}

func (c *Comment) EscapeContent() {
	c.Content = html.EscapeString(c.Content)
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByCommentID(ctx context.Context, commentID int64) (*Comment, bool, error)
	ListNSubComments(ctx context.Context, rootID []int64) ([]*Comment, error)
	GetLikeCount(ctx context.Context, commentID int64) (int64, bool, error)
	UpdateLikeCount(ctx context.Context, commentID int64, count int64) error
	List(ctx context.Context, commentQuery *CommentQuery) (paginator.Page[*Comment], error)
}
