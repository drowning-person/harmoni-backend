package events

import (
	"context"
	eventlike "harmoni/app/harmoni/api/mq/v1/like"
	"harmoni/app/notification/internal/entity/remind"
	usercaseremind "harmoni/app/notification/internal/usecase/remind"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"
)

type RemindEventsHandler struct {
	ru *usercaseremind.RemindUsecase
}

func NewLikeEventsHandler(
	ru *usercaseremind.RemindUsecase,
) *RemindEventsHandler {
	return &RemindEventsHandler{
		ru: ru,
	}
}

func fromLikeType(t eventlike.LikeType) object.ObjectType {
	switch t {
	case eventlike.LikeType_LikePost:
		return object.ObjectTypePost
	case eventlike.LikeType_LikeComment:
		return object.ObjectTypeComment
	default:
		return 0
	}
}

func (h *RemindEventsHandler) HandleLikeCreated(ctx context.Context, msg *eventlike.LikeCreatedMessage) error {
	if msg.IsCancel || msg.BaseMessage.LikeType == eventlike.LikeType_LikeUser ||
		msg.BaseMessage.LikeType == eventlike.LikeType_LikeNo {
		return nil
	}

	createdAt := msg.CreatedAt.AsTime()
	return h.ru.Create(ctx, &remind.CreateReq{
		RecipientID:  msg.TargetUserID,
		SenderIDs:    []int64{msg.UserID},
		Action:       action.ActionLike,
		ObjectID:     msg.LikingID,
		ObjectType:   fromLikeType(msg.BaseMessage.LikeType),
		LastReadTime: &createdAt,
	})
}
