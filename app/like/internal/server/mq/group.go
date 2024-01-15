package mq

import (
	"encoding/json"

	v1 "harmoni/api/like/mq/v1"
	eventlike "harmoni/app/like/internal/usecase/like/events"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/mq"
	"harmoni/internal/pkg/mq/subscriber"
	"harmoni/internal/pkg/watermillkratos"
	"harmoni/internal/types/events"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	groupName = "like"
)

func NewLikeGroup(
	conf *conf.MessageQueue,
	r *message.Router,
	remindEventsHandler *eventlike.LikeEventsHandler,
	logger *log.Helper,
	m ...message.HandlerMiddleware,
) error {
	sub, err := subscriber.NewSubscriber((&mq.MessageQueue{}).FromConfig(conf), groupName, watermillkratos.NewLogger(logger, "msg"))
	if err != nil {
		return err
	}
	g := &mq.Group{
		Router: r,
		Name:   groupName,
		Sub:    sub,
	}
	g.Handle(events.TopicLikeCreated, func(msg *message.Message) error {
		var m v1.LikeCreatedMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return remindEventsHandler.HandleLikeCreated(msg.Context(), &m)
	})
	return nil
}
