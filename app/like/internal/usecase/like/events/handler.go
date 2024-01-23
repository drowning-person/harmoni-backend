package events

import (
	"context"

	objectv1 "harmoni/api/common/object/v1"
	v1 "harmoni/api/like/mq/v1"
	apiuser "harmoni/app/harmoni/api/grpc/v1/user"
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/internal/types/iface"

	"github.com/go-kratos/kratos/v2/log"
)

type LikeEventsHandler struct {
	likeRepo entitylike.LikeRepository
	txMgr    iface.Transaction
	logger   *log.Helper
}

func NewLikeEventsHandler(
	likeRepo entitylike.LikeRepository,
	txMgr iface.Transaction,
	logger log.Logger,
) *LikeEventsHandler {
	return &LikeEventsHandler{
		likeRepo: likeRepo,
		txMgr:    txMgr,
		logger:   log.NewHelper(log.With(logger, "module", "events/like", "service", "like")),
	}
}

func FromEventLikeType(likeType v1.LikeType) entitylike.LikeType {
	switch likeType {
	case v1.LikeType_LikePost:
		return entitylike.LikePost
	case v1.LikeType_LikeComment:
		return entitylike.LikeComment
	case v1.LikeType_LikeUser:
		return entitylike.LikeUser
	}
	return entitylike.LikePost
}

func (h *LikeEventsHandler) HandleLikeCreated(ctx context.Context, msg *v1.LikeCreatedMessage) error {
	return h.txMgr.ExecTx(ctx, func(ctx context.Context) error {
		err := h.likeRepo.Save(ctx, &entitylike.Like{
			LikingID:   msg.LikingID,
			LikeType:   FromEventLikeType(msg.BaseMessage.LikeType),
			User:       &apiuser.UserBasic{Id: msg.UserID},
			TargetUser: &apiuser.UserBasic{Id: msg.TargetUserID},
			ObjectID:   msg.ObjectID,
		}, msg.IsCancel)
		if err != nil {
			return err
		}
		return h.likeRepo.AddLikeCount(ctx, &objectv1.Object{
			Id:   msg.TargetUserID,
			Type: objectv1.ObjectType_OBJECT_TYPE_USER,
		}, 1)
	})
}
