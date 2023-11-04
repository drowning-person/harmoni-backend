package events

import (
	"context"
	userentity "harmoni/internal/entity/user"
	event "harmoni/internal/types/events/like"
)

type UserEventsHandler struct {
	userRepo userentity.UserRepository
}

func NewUserEventsHandler(userRepo userentity.UserRepository) *UserEventsHandler {
	return &UserEventsHandler{
		userRepo: userRepo,
	}
}

func (h *UserEventsHandler) HandleLikeStore(ctx context.Context, msg *event.LikeStoreMessage) error {
	if msg.LikeType != event.LikeUser {
		return nil
	}
	for k, v := range msg.Counts {
		if err := h.userRepo.UpdateLikeCount(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}
