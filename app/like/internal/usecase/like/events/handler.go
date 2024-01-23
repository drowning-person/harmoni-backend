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

func (h *LikeEventsHandler) HandleLikeCreated(ctx context.Context, msg *v1.LikeCreatedMessage) error {
	return h.txMgr.ExecTx(ctx, func(ctx context.Context) error {
		err := h.likeRepo.Save(ctx, &entitylike.Like{
			LikingID:   msg.GetLikingId(),
			ObjectType: msg.BaseMessage.GetObjectType(),
			User:       &apiuser.UserBasic{Id: msg.GetUserId()},
			TargetUser: &apiuser.UserBasic{Id: msg.GetTargetUserId()},
			ObjectID:   msg.GetObjectId(),
		}, msg.IsCancel)
		if err != nil {
			return err
		}
		if entitylike.ShouldAddUserLikeCount(msg.BaseMessage.GetObjectType()) {
			return h.likeRepo.AddLikeCount(ctx, &objectv1.Object{
				Id:   msg.GetTargetUserId(),
				Type: objectv1.ObjectType_OBJECT_TYPE_USER,
			}, 1)
		}
		return nil
	})
}
