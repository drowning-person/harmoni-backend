package like

import (
	"encoding/json"
	"harmoni/internal/infrastructure/mq"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/like/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewLikeGroup(
	r *message.Router,
	likeEventsHandler *events.LikeEventsHandler,
	sub message.Subscriber,
	m ...message.HandlerMiddleware,
) {
	g := &mq.Group{
		Router: r,
		Name:   "like",
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeCreated, func(msg *message.Message) error {
		var m eventlike.LikeCreatedMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return likeEventsHandler.HandleLikeCreated(msg.Context(), &m)
	})
}
