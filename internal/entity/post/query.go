package post

import "harmoni/internal/entity"

const (
	PostOrderByCreatedTime = "ctime"
	PostOrderByReplyTime   = "rtime"
	PostOrderByLike        = "like"
)

type PostQuery struct {
	entity.PageCond
	// query condition
	QueryCond string
	// tag condition
	TagID int64
	// author condition
	AuthorIDs []int64
	// user id
	UserID int64
}
