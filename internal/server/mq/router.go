package mq

import (
	"context"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/infrastructure/mq/subscriber"
	"harmoni/internal/server/mq/group/comment"
	"harmoni/internal/server/mq/group/like"
	"harmoni/internal/server/mq/group/post"
	"harmoni/internal/server/mq/group/user"
	"harmoni/internal/types/iface"
	commentevent "harmoni/internal/usecase/comment/events"
	likeevent "harmoni/internal/usecase/like/events"
	postevent "harmoni/internal/usecase/post/events"
	userevent "harmoni/internal/usecase/user/events"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/wire"
)

var _ iface.Executor = (*MQExecutor)(nil)

type MQExecutor struct {
	*message.Router
}

var ProviderSetMQ = wire.NewSet(
	NewMQRouter,
	NewExecutor,
)

func (r *MQExecutor) Start() error {
	return r.Run(context.Background())
}

func (r *MQExecutor) Shutdown() error {
	return r.Close()
}

func NewExecutor(r *message.Router) *MQExecutor {
	return &MQExecutor{
		Router: r,
	}
}

func NewMQRouter(
	conf *config.MessageQueue,
	commentEventsHandler *commentevent.CommentEventsHandler,
	likeevent *likeevent.LikeEventsHandler,
	postEventsHandler *postevent.PostEventsHandler,
	userevent *userevent.UserEventsHandler,
) (*message.Router, error) {
	r, err := message.NewRouter(message.RouterConfig{}, watermill.NewStdLogger(true, true))
	if err != nil {
		return nil, err
	}
	var sub message.Subscriber
	if conf.RabbitMQ != nil {
		sub, err = subscriber.NewAMQPSubscriber(conf.RabbitMQ)
		if err != nil {
			return nil, err
		}
	}

	comment.NewCommentGroup(r, commentEventsHandler, sub)
	like.NewLikeGroup(r, likeevent, sub)
	post.NewPostGroup(r, postEventsHandler, sub)
	user.NewUserGroup(r, userevent, sub)
	return r, nil
}
