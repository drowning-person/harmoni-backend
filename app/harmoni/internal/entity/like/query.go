package like

import "harmoni/app/harmoni/internal/entity"

type LikeQuery struct {
	entity.PageCond
	// object id
	ObjectID int64
	// user id
	UserID int64
	Type   LikeType
}
