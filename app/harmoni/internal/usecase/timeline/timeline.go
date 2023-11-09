package timeline

import (
	"context"
	followentity "harmoni/app/harmoni/internal/entity/follow"
	"harmoni/app/harmoni/internal/entity/paginator"
	postentity "harmoni/app/harmoni/internal/entity/post"
	"harmoni/app/harmoni/internal/usecase/post"

	"go.uber.org/zap"
)

type TimeLinePullUsecase struct {
	followRepo  followentity.FollowRepository
	postUseCase *post.PostUseCase
	logger      *zap.SugaredLogger
}

func NewTimeLineUsecase(
	followRepo followentity.FollowRepository,
	postUseCase *post.PostUseCase,
	loggger *zap.SugaredLogger,
) *TimeLinePullUsecase {
	return &TimeLinePullUsecase{
		followRepo:  followRepo,
		postUseCase: postUseCase,
		logger:      loggger,
	}
}

// GetTimeLineByUserID Get user's timeline
func (u *TimeLinePullUsecase) GetTimeLineByUserID(ctx context.Context, authorID int64, userID int64, queryCond *postentity.PostQuery) (*paginator.Page[postentity.PostBasicInfo], error) {
	return u.postUseCase.GetPage(ctx, &postentity.PostQuery{
		PageCond:  queryCond.PageCond,
		AuthorIDs: []int64{authorID},
		UserID:    userID,
	})
}

// GetTimeLine get user's followings timeline
func (u *TimeLinePullUsecase) GetTimeLine(ctx context.Context, userID int64, queryCond *postentity.PostQuery) (*paginator.Page[postentity.PostBasicInfo], error) {
	followings, err := u.followRepo.GetFollowingUsersAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.postUseCase.GetPage(ctx, &postentity.PostQuery{
		PageCond:  queryCond.PageCond,
		AuthorIDs: followings,
		UserID:    userID,
	})
}
