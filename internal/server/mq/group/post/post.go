package post

import (
	"encoding/json"
	"harmoni/internal/infrastructure/mq"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/post/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewPostGroup(
	r *message.Router,
	postEventsHandler *events.PostEventsHandler,
	sub message.Subscriber,
	m ...message.HandlerMiddleware,
) {
	g := &mq.Group{
		Router: r,
		Name:   "post",
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeStore, func(msg *message.Message) error {
		var m eventlike.LikeStoreMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return postEventsHandler.HandleLikeStore(msg.Context(), &m)
	})
}
