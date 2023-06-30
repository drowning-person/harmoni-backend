package usecase

import (
	"context"
	followentity "harmoni/internal/entity/follow"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"

	"go.uber.org/zap"
)

type TimeLinePullUsecase struct {
	followRepo followentity.FollowRepository
	postRepo   postentity.PostRepository
	logger     *zap.SugaredLogger
}

func NewTimeLineUsecase(
	followRepo followentity.FollowRepository,
	postRepo postentity.PostRepository,
	loggger *zap.SugaredLogger,
) *TimeLinePullUsecase {
	return &TimeLinePullUsecase{
		followRepo: followRepo,
		postRepo:   postRepo,
		logger:     loggger,
	}
}

func (u *TimeLinePullUsecase) GetTimeLineByUserID(ctx context.Context, userID int64, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	return u.postRepo.GetByUserID(ctx, userID, queryCond)
}

func (u *TimeLinePullUsecase) GetTimeLine(ctx context.Context, userID int64, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	followings, err := u.followRepo.GetFollowingUsersAll(ctx, userID)
	if err != nil {
		return paginator.Page[postentity.Post]{}, err
	}

	return u.postRepo.GetByUserIDs(ctx, followings, queryCond)
}
