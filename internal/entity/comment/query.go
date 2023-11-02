package comment

import "harmoni/internal/entity"

type CommentQuery struct {
	entity.PageCond
	// object id
	ObjectID string
	// root id
	RootID int64
	// query condition
	QueryCond string
	// user id
	UserID int64
}
