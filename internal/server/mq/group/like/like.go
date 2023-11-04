package like

import (
	"encoding/json"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/infrastructure/mq"
	"harmoni/internal/infrastructure/mq/subscriber"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/like/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	groupName = "like"
)

func NewLikeGroup(
	conf *config.MessageQueue,
	r *message.Router,
	likeEventsHandler *events.LikeEventsHandler,
	m ...message.HandlerMiddleware,
) error {
	sub, err := subscriber.NewSubscriber(conf, groupName)
	if err != nil {
		return err
	}
	g := &mq.Group{
		Router: r,
		Name:   groupName,
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeCreated, func(msg *message.Message) error {
		var m eventlike.LikeCreatedMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return likeEventsHandler.HandleLikeCreated(msg.Context(), &m)
	})
	return nil
}
