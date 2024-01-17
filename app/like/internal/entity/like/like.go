package like

import (
	"context"

	mqlike "harmoni/api/like/mq/v1"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/internal/pkg/paginator"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type LikeType uint8

const (
	LikePost LikeType = iota + 1
	LikeComment
	LikeUser
)

type Like struct {
	LikingID   int64
	User       *v1.UserBasic
	LikeType   LikeType
	TargetUser *v1.UserBasic
	ObjectID   int64
}

func (l *Like) ToLikeCreateMessage(isCancel bool) *mqlike.LikeCreatedMessage {
	return &mqlike.LikeCreatedMessage{
		BaseMessage: &mqlike.BaseMessage{
			LikeType:  mqlike.LikeType(l.LikeType),
			CreatedAt: timestamppb.Now(),
		},
		LikingID:     l.LikingID,
		UserID:       l.User.GetId(),
		TargetUserID: l.TargetUser.GetId(),
		ObjectID:     l.ObjectID,
		IsCancel:     isCancel,
	}
}

var (
	LikeTypeList = []LikeType{LikeUser, LikePost, LikeComment}
)

type ListLikeObjectQuery struct {
	paginator.PageRequest
	UserID   int64
	LikeType LikeType
}

type LikeRepository interface {
	Save(ctx context.Context, like *Like, isCancel bool) error
	Get(ctx context.Context, like *Like) (*Like, error)
	IsExist(ctx context.Context, like *Like) (bool, error)
	ListLikeObjectByUserID(ctx context.Context, query *ListLikeObjectQuery) ([]*Like, int64, error)
}
