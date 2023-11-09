package events

import (
	"context"
	likeentity "harmoni/app/harmoni/internal/entity/like"
	event "harmoni/app/harmoni/internal/types/events/like"
)

type LikeEventsHandler struct {
	likeRepo likeentity.LikeRepository
}

func NewLikeEventsHandler(likeRepo likeentity.LikeRepository) *LikeEventsHandler {
	return &LikeEventsHandler{
		likeRepo: likeRepo,
	}
}

func FromEventLikeType(likeType event.LikeType) likeentity.LikeType {
	switch likeType {
	case event.LikePost:
		return likeentity.LikePost
	case event.LikeComment:
		return likeentity.LikeComment
	case event.LikeUser:
		return likeentity.LikeUser
	}
	return likeentity.LikePost
}

func (h *LikeEventsHandler) HandleLikeCreated(ctx context.Context, msg *event.LikeCreatedMessage) error {
	return h.likeRepo.Save(ctx, &likeentity.Like{
		UserID:       msg.UserID,
		TargetUserID: msg.TargetUserID,
		LikingID:     msg.LikingID,
		LikeType:     FromEventLikeType(msg.LikeType),
		Canceled:     msg.IsCancel,
	})
}
