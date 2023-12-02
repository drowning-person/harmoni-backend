package mq

import (
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/server/mq/group/comment"
	"harmoni/app/harmoni/internal/server/mq/group/like"
	"harmoni/app/harmoni/internal/server/mq/group/post"
	"harmoni/app/harmoni/internal/server/mq/group/user"
	commentevent "harmoni/app/harmoni/internal/usecase/comment/events"
	likeevent "harmoni/app/harmoni/internal/usecase/like/events"
	postevent "harmoni/app/harmoni/internal/usecase/post/events"
	userevent "harmoni/app/harmoni/internal/usecase/user/events"
	"harmoni/internal/pkg/server"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/garsue/watermillzap"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSetMQ = wire.NewSet(
	NewMQRouter,
	server.NewMQServer,
)

func NewMQRouter(
	conf *config.MessageQueue,
	commentEventsHandler *commentevent.CommentEventsHandler,
	likeevent *likeevent.LikeEventsHandler,
	postEventsHandler *postevent.PostEventsHandler,
	userevent *userevent.UserEventsHandler,
	logger *zap.Logger,
) (*message.Router, error) {
	r, err := message.NewRouter(message.RouterConfig{}, watermillzap.NewLogger(logger))
	if err != nil {
		return nil, err
	}
	r.AddMiddleware(middleware.Retry{
		MaxRetries: 3,
	}.Middleware)
	{
		err = comment.NewCommentGroup(conf, r, commentEventsHandler, logger)
		if err != nil {
			return nil, err
		}
	}
	{
		err = like.NewLikeGroup(conf, r, likeevent, logger)
		if err != nil {
			return nil, err
		}
	}
	{
		err = post.NewPostGroup(conf, r, postEventsHandler, logger)
		if err != nil {
			return nil, err
		}
	}
	{
		err = user.NewUserGroup(conf, r, userevent, logger)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}
