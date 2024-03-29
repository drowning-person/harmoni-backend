package events

import (
	"context"
	v1 "harmoni/app/harmoni/api/mq/v1/like"
	"harmoni/app/harmoni/internal/entity"
	likeentity "harmoni/app/harmoni/internal/entity/like"
)

type LikeEventsHandler struct {
	likeRepo likeentity.LikeRepository
}

func NewLikeEventsHandler(likeRepo likeentity.LikeRepository) *LikeEventsHandler {
	return &LikeEventsHandler{
		likeRepo: likeRepo,
	}
}

func FromEventLikeType(likeType v1.LikeType) likeentity.LikeType {
	switch likeType {
	case v1.LikeType_LikePost:
		return likeentity.LikePost
	case v1.LikeType_LikeComment:
		return likeentity.LikeComment
	case v1.LikeType_LikeUser:
		return likeentity.LikeUser
	}
	return likeentity.LikePost
}

func (h *LikeEventsHandler) HandleLikeCreated(ctx context.Context, msg *v1.LikeCreatedMessage) error {
	return h.likeRepo.Save(ctx, &likeentity.Like{
		UserID:       msg.UserID,
		TargetUserID: msg.TargetUserID,
		LikingID:     msg.LikingID,
		LikeType:     FromEventLikeType(msg.BaseMessage.LikeType),
		Canceled:     msg.IsCancel,
		TimeMixin: entity.TimeMixin{
			CreatedAt: msg.CreatedAt.AsTime(),
		},
	})
}
