package comment

import "harmoni/internal/entity"

type CommentQuery struct {
	entity.PageCond
	// object id
	ObjectID string
	// root id
	RootID string
	// query condition
	QueryCond string
	// user id
	UserID int64
}
