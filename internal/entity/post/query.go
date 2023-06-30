package post

import "harmoni/internal/entity"

type PostQuery struct {
	entity.PageCond
	// query condition
	QueryCond string
	// user id
	UserID int64
}
