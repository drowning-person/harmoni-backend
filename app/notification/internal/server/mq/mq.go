package mq

import (
	"harmoni/app/notification/internal/server/mq/group/notification"
	"harmoni/app/notification/internal/usecase/remind/events"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/server"
	"harmoni/internal/pkg/watermillkratos"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSetMQ = wire.NewSet(
	NewMQRouter,
	server.NewMQServer,
)

func NewMQRouter(
	conf *conf.MessageQueue,
	remindEventsHandler *events.RemindEventsHandler,
	logger log.Logger,
) (*message.Router, error) {
	helper := log.NewHelper(log.With(logger, "module", "message-queue"))
	r, err := message.NewRouter(message.RouterConfig{}, watermillkratos.NewLogger(helper, log.DefaultMessageKey))
	if err != nil {
		return nil, err
	}
	r.AddMiddleware(middleware.Retry{
		MaxRetries: 3,
	}.Middleware)
	{
		err = notification.NewNotificationGroup(conf, r, remindEventsHandler, helper)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
