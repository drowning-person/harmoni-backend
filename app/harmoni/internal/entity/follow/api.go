package follow

import (
	"harmoni/app/harmoni/internal/entity"
	"harmoni/app/harmoni/internal/entity/paginator"
	userentity "harmoni/app/harmoni/internal/entity/user"
)

type FollowRequest struct {
	Type     FollowedType `json:"type,omitempty" validate:"required"`
	UserID   int64        `json:"user_id,omitempty"`
	ObjectID int64        `json:"oid,omitempty,string" validate:"required"`
	IsCancel bool         `json:"is_cancel,omitempty"`
}

type FollowReply struct {
}

type GetFollowingsRequest struct {
	entity.PageCond
	Type   FollowedType `query:"type,omitempty" validate:"required"`
	UserID int64        `json:"user_id,omitempty"`
}

type GetFollowingsReply[T any] struct {
	paginator.Page[T]
}

type GetFollowersRequest struct {
	entity.PageCond
	UserID int64 `json:"user_id,omitempty"`
}

type GetFollowerReply struct {
	paginator.Page[userentity.UserBasicInfo]
}

type UnFollowRequest struct {
	Type     FollowedType `json:"type,omitempty" validate:"required"`
	UserID   int64        `json:"user_id,omitempty"`
	ObjectID int64        `json:"oid,omitempty,string" validate:"required"`
}

type UnFollowReply struct {
}

type IsFollowingRequest struct {
	Type     FollowedType `query:"type,omitempty" validate:"required"`
	UserID   int64        `json:"user_id,omitempty"`
	ObjectID int64        `query:"oid,omitempty,string" validate:"required"`
}

type IsFollowingReply struct {
	Following bool `json:"following"`
}

type AreFollowEachOtherRequest struct {
	UserID   int64 `json:"user_id,omitempty"`
	ObjectID int64 `query:"oid,omitempty,string" validate:"required"`
}

type AreFollowEachOtherReply struct {
	FollowedEachOther bool `json:"followed_each_other"`
}
