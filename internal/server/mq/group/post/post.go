package post

import (
	"encoding/json"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/infrastructure/mq"
	"harmoni/internal/infrastructure/mq/subscriber"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/post/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	groupName = "post"
)

func NewPostGroup(
	conf *config.MessageQueue,
	r *message.Router,
	postEventsHandler *events.PostEventsHandler,
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
	g.Handle(eventlike.TopicLikeStore, func(msg *message.Message) error {
		var m eventlike.LikeStoreMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return postEventsHandler.HandleLikeStore(msg.Context(), &m)
	})
	return nil
}
