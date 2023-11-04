package events

import (
	"context"
	postentity "harmoni/internal/entity/post"
	event "harmoni/internal/types/events/like"
)

type PostEventsHandler struct {
	postRepo postentity.PostRepository
}

func NewPostEventsHandler(postRepo postentity.PostRepository) *PostEventsHandler {
	return &PostEventsHandler{
		postRepo: postRepo,
	}
}

func (h *PostEventsHandler) HandleLikeStore(ctx context.Context, msg *event.LikeStoreMessage) error {
	if msg.LikeType != event.LikePost {
		return nil
	}
	for k, v := range msg.Counts {
		if err := h.postRepo.UpdateLikeCount(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}
