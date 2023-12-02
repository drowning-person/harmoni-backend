package notification

import (
	"encoding/json"
	eventlike "harmoni/app/harmoni/api/mq/v1/like"
	eventremind "harmoni/app/notification/internal/usecase/remind/events"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/mq"
	"harmoni/internal/pkg/mq/subscriber"
	"harmoni/internal/pkg/watermillkratos"
	"harmoni/internal/types/events"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	groupName = "notification"
)

func NewNotificationGroup(
	conf *conf.MessageQueue,
	r *message.Router,
	remindEventsHandler *eventremind.RemindEventsHandler,
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
		var m eventlike.LikeCreatedMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return remindEventsHandler.HandleLikeCreated(msg.Context(), &m)
	})
	return nil
}
