package like

import (
	"context"

	v1 "harmoni/app/harmoni/api/grpc/v1/user"
)

type LikeType uint8

const (
	LikePost LikeType = iota + 1
	LikeComment
	LikeUser
)

type Like struct {
	LikingID       int64
	User           *v1.UserBasic
	LikeType       LikeType
	TargetUser     *v1.UserBasic
	TargetObjectID int64
}

var (
	LikeTypeList = []LikeType{LikeUser, LikePost, LikeComment}
)

type LikeRepository interface {
	Save(ctx context.Context, like *Like, isCancel bool) error
}
