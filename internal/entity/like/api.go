package like

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
)

type LikeRequest struct {
	Type     LikeType `json:"type,omitempty" validate:"required"`
	UserID   int64    `json:"user_id,omitempty"`
	ObjectID int64    `json:"oid,omitempty,string" validate:"required"`
	IsCancel bool     `json:"is_cancel,omitempty"`
}

type LikeReply struct {
}

type GetLikingsRequest struct {
	entity.PageCond
	Type         LikeType `query:"type,omitempty" validate:"required"`
	UserID       int64    `json:"-"`
	TargetUserID int64    `query:"user_id" validate:"required"`
}

type GetLikingsReply[T any] struct {
	paginator.Page[T]
}

type IsLikingRequest struct {
	Type     LikeType `query:"type,omitempty" validate:"required"`
	UserID   int64    `json:"user_id,omitempty"`
	ObjectID int64    `query:"oid,omitempty,string" validate:"required"`
}

type IsLikingReply struct {
	Liking bool `json:"liking"`
}
