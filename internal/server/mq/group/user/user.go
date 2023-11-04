package user

import (
	"encoding/json"
	"harmoni/internal/infrastructure/mq"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/user/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewUserGroup(
	r *message.Router,
	userEventsHandler *events.UserEventsHandler,
	sub message.Subscriber,
	m ...message.HandlerMiddleware,
) {
	g := &mq.Group{
		Router: r,
		Name:   "user",
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeStore, func(msg *message.Message) error {
		var m eventlike.LikeStoreMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return userEventsHandler.HandleLikeStore(msg.Context(), &m)
	})
}
