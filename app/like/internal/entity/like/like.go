package like

import (
	"context"

	objectv1 "harmoni/api/common/object/v1"
	mqlike "harmoni/api/like/mq/v1"
	userv1 "harmoni/api/user/v1"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/like/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/paginator"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Like struct {
	LikingID   int64
	User       *v1.UserBasic
	TargetUser *v1.UserBasic
	ObjectID   int64
	ObjectType objectv1.ObjectType
}

func (l *Like) ToLikeCreateMessage(isCancel bool) *mqlike.LikeCreatedMessage {
	return &mqlike.LikeCreatedMessage{
		BaseMessage: &mqlike.BaseMessage{
			ObjectType: l.ObjectType,
			CreatedAt:  timestamppb.Now(),
		},
		LikingId:     l.LikingID,
		UserId:       l.User.GetId(),
		TargetUserId: l.TargetUser.GetId(),
		ObjectId:     l.ObjectID,
		IsCancel:     isCancel,
	}
}

type QueryType int8

const (
	QueryTypeUser QueryType = iota + 1
	QueryTypeObject
)

type ListLikeObjectQuery struct {
	paginator.PageRequest
	UserID     int64
	ObjectType objectv1.ObjectType
}

type ListObjectLikedUserQuery struct {
	paginator.PageRequest
	ObjectID   int64
	ObjectType objectv1.ObjectType
}

func ShouldAddUserLikeCount(objectType objectv1.ObjectType) bool {
	return objectType == objectv1.ObjectType_OBJECT_TYPE_POST
}

func (l *Like) Validate() error {
	if l.User == nil || l.TargetUser == nil {
		return errorx.NotFound(userv1.UserNotFound)
	}
	if l.User.GetId() == l.TargetUser.GetId() {
		return errorx.BadRequest(reason.DisallowLikeYourSelf)
	}
	if l.ObjectType != objectv1.ObjectType_OBJECT_TYPE_POST &&
		l.ObjectType != objectv1.ObjectType_OBJECT_TYPE_COMMENT {
		return errorx.BadRequest(reason.LikeUnknownType)
	}
	return nil
}

type LikeRepository interface {
	Save(ctx context.Context, like *Like, isCancel bool) error
	Get(ctx context.Context, like *Like) (*Like, error)
	IsExist(ctx context.Context, like *Like) (bool, error)
	ListLikeObjectByUserID(ctx context.Context, query *ListLikeObjectQuery) ([]*Like, int64, error)
	ListObjectLikedUserByObjectID(ctx context.Context, query *ListObjectLikedUserQuery) ([]*Like, int64, error)

	// count
	ObjectLikeCount(ctx context.Context, object *objectv1.Object) (*LikeCount, error)
	ListObjectLikeCount(ctx context.Context, objectIDs []int64, objectType objectv1.ObjectType) (LikeCountList, error)
	AddLikeCount(ctx context.Context, object *objectv1.Object, count int64) error
}
