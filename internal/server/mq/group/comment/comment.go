package comment

import (
	"encoding/json"
	"harmoni/internal/infrastructure/mq"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/usecase/comment/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewCommentGroup(
	r *message.Router,
	commentEventsHandler *events.CommentEventsHandler,
	sub message.Subscriber,
	m ...message.HandlerMiddleware,
) {
	g := &mq.Group{
		Router: r,
		Name:   "comment",
		Sub:    sub,
	}
	g.Handle(eventlike.TopicLikeStore, func(msg *message.Message) error {
		var m eventlike.LikeStoreMessage
		if err := json.Unmarshal(msg.Payload, &m); err != nil {
			return err
		}
		return commentEventsHandler.HandleLikeStore(msg.Context(), &m)
	})
}
