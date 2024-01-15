package mq

import (
	"harmoni/app/like/internal/usecase/like/events"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/mq"
	"harmoni/internal/pkg/mq/publisher"
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
	likeEventsHandler *events.LikeEventsHandler,
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
		err = NewLikeGroup(conf, r, likeEventsHandler, helper)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func NewPublisher(
	conf *conf.MessageQueue,
	logger log.Logger,
) (message.Publisher, error) {
	return publisher.NewPublisher((&mq.MessageQueue{}).FromConfig(conf), logger)
}
