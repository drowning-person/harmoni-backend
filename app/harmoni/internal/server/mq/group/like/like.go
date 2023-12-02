package like

import (
	"encoding/json"
	v1 "harmoni/app/harmoni/api/mq/v1/like"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/infrastructure/mq"
	"harmoni/app/harmoni/internal/infrastructure/mq/subscriber"
	eventlike "harmoni/app/harmoni/internal/types/events/like"
	"harmoni/app/harmoni/internal/usecase/like/events"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

const (
	groupName = "like"
)

func NewLikeGroup(
	conf *config.MessageQueue,
	r *message.Router,
	likeEventsHandler *events.LikeEventsHandler,
	logger *zap.Logger,
	m ...message.HandlerMiddleware,
) error {
	sub, err := subscriber.NewSubscriber(conf, groupName, logger)
	if err != nil {
		return err
	}
	g := &mq.Group{
		Router: r,
		Name:   groupName,
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeCreated, func(msg *message.Message) error {
		var m v1.LikeCreatedMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return likeEventsHandler.HandleLikeCreated(msg.Context(), &m)
	})
	return nil
}
