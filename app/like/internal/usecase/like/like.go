package like

import (
	"context"

	objectv1 "harmoni/api/common/object/v1"
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
	like := &entitylike.Like{
		User:       &v1.UserBasic{Id: req.UserID},
		ObjectType: req.ObjectType,
		TargetUser: &v1.UserBasic{Id: req.TargetUserID},
		ObjectID:   req.TargetObjectID,
	}
	err := like.Validate()
	if err != nil {
		return err
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
	list, total, err := u.likeRepo.ListLikeObjectByUserID(ctx, query)
	if err != nil {
		return nil, 0, err
	} else if total == 0 {
		return nil, 0, errorx.NotFound(reason.LikeNotExist)
	}
	return list, total, nil
}

func (u *LikeUsecase) ListLikeUserByObjectID(ctx context.Context, query *entitylike.ListObjectLikedUserQuery) ([]*entitylike.Like, int64, error) {
	return u.likeRepo.ListObjectLikedUserByObjectID(ctx, query)
}

func (u *LikeUsecase) LikeCount(ctx context.Context, object *objectv1.Object) (*entitylike.LikeCount, error) {
	return u.likeRepo.ObjectLikeCount(ctx, object)
}

func (u *LikeUsecase) ListLikeCount(ctx context.Context, objectIDs []int64, objectType objectv1.ObjectType) (entitylike.LikeCountList, error) {
	return u.likeRepo.ListObjectLikeCount(ctx, objectIDs, objectType)
}
