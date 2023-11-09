package user

import (
	"encoding/json"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/infrastructure/mq"
	"harmoni/app/harmoni/internal/infrastructure/mq/subscriber"
	eventlike "harmoni/app/harmoni/internal/types/events/like"
	"harmoni/app/harmoni/internal/usecase/user/events"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

var (
	groupName = "user"
)

func NewUserGroup(
	conf *config.MessageQueue,
	r *message.Router,
	userEventsHandler *events.UserEventsHandler,
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
	g.Handle(eventlike.TopicLikeStore, func(msg *message.Message) error {
		var m eventlike.LikeStoreMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return userEventsHandler.HandleLikeStore(msg.Context(), &m)
	})
	return nil
}
