package events

import (
	"context"

	v1 "harmoni/api/like/mq/v1"
	apiuser "harmoni/app/harmoni/api/grpc/v1/user"
	entitylike "harmoni/app/like/internal/entity/like"

	"github.com/go-kratos/kratos/v2/log"
)

type LikeEventsHandler struct {
	likeRepo entitylike.LikeRepository
	logger   *log.Helper
}

func NewLikeEventsHandler(
	likeRepo entitylike.LikeRepository,
	logger log.Logger,
) *LikeEventsHandler {
	return &LikeEventsHandler{
		likeRepo: likeRepo,
		logger:   log.NewHelper(log.With(logger, "module", "events/like", "service", "like")),
	}
}

func FromEventLikeType(likeType v1.LikeType) entitylike.LikeType {
	switch likeType {
	case v1.LikeType_LikePost:
		return entitylike.LikePost
	case v1.LikeType_LikeComment:
		return entitylike.LikeComment
	case v1.LikeType_LikeUser:
		return entitylike.LikeUser
	}
	return entitylike.LikePost
}

func (h *LikeEventsHandler) HandleLikeCreated(ctx context.Context, msg *v1.LikeCreatedMessage) error {
	return h.likeRepo.Save(ctx, &entitylike.Like{
		LikingID:   msg.LikingID,
		LikeType:   FromEventLikeType(msg.BaseMessage.LikeType),
		User:       &apiuser.UserBasic{Id: msg.UserID},
		TargetUser: &apiuser.UserBasic{Id: msg.TargetUserID},
		ObjectID:   msg.ObjectID,
	}, msg.IsCancel)
}
