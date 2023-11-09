package events

import (
	"context"
	userentity "harmoni/app/harmoni/internal/entity/comment"
	event "harmoni/app/harmoni/internal/types/events/like"
)

type CommentEventsHandler struct {
	userRepo userentity.CommentRepository
}

func NewCommentEventsHandler(userRepo userentity.CommentRepository) *CommentEventsHandler {
	return &CommentEventsHandler{
		userRepo: userRepo,
	}
}

func (h *CommentEventsHandler) HandleLikeStore(ctx context.Context, msg *event.LikeStoreMessage) error {
	if msg.LikeType != event.LikeComment {
		return nil
	}
	for k, v := range msg.Counts {
		if err := h.userRepo.UpdateLikeCount(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}
