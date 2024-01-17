package like

import (
	"context"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/app/like/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/types/events"
	"harmoni/internal/types/iface"

	"github.com/go-kratos/kratos/v2/log"
)

type LikeUsecase struct {
	likeRepo  entitylike.LikeRepository
	publisher iface.Publisher
	logger    *log.Helper
}

func NewLikeUsecase(
	likeRepo entitylike.LikeRepository,
	publisher iface.Publisher,
	logger log.Logger,
) *LikeUsecase {
	return &LikeUsecase{
		likeRepo:  likeRepo,
		publisher: publisher,
		logger:    log.NewHelper(log.With(logger, "module", "usecase/like", "service", "like")),
	}
}

func (u *LikeUsecase) Like(ctx context.Context, req *LikeRequest) error {
	if req.UserID == req.TargetUserID {
		return errorx.BadRequest(reason.DisallowLikeYourSelf)
	}
	like := &entitylike.Like{
		User:       &v1.UserBasic{Id: req.UserID},
		LikeType:   req.LikeType,
		TargetUser: &v1.UserBasic{Id: req.TargetUserID},
		ObjectID:   req.TargetObjectID,
	}
	exist, err := u.likeRepo.IsExist(ctx, like)
	if err != nil {
		return err
	}

	if req.IsCancel && !exist {
		return errorx.NotFound(reason.LikeNotExist)
	} else if !req.IsCancel && exist {
		return errorx.BadRequest(reason.LikeAlreadyExist)
	}

	return u.publisher.Publish(ctx, events.TopicLikeCreated, like.ToLikeCreateMessage(req.IsCancel))
}

func (u *LikeUsecase) ListLikeObjectByUserID(ctx context.Context, query *entitylike.ListLikeObjectQuery) ([]*entitylike.Like, int64, error) {
	return u.likeRepo.ListLikeObjectByUserID(ctx, query)
}
