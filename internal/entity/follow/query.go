package follow

import "harmoni/internal/entity"

type FollowQuery struct {
	entity.PageCond
	// object id
	ObjectID int64
	// user id
	UserID int64
	Type   FollowedType
}
