package follow

import (
	"context"
	followentity "harmoni/app/harmoni/internal/entity/follow"
	"harmoni/app/harmoni/internal/entity/paginator"
	"harmoni/app/harmoni/internal/pkg/reason"
	taguse "harmoni/app/harmoni/internal/usecase/tag"
	useruse "harmoni/app/harmoni/internal/usecase/user"
	"harmoni/internal/pkg/errorx"

	"go.uber.org/zap"
)

type FollowUseCase struct {
	followRepo  followentity.FollowRepository
	userUseCase *useruse.UserUseCase
	tagUseCase  *taguse.TagUseCase
	logger      *zap.SugaredLogger
}

func NewFollowUseCase(followRepo followentity.FollowRepository,
	userUseCase *useruse.UserUseCase,
	tagUseCase *taguse.TagUseCase,
	logger *zap.SugaredLogger) *FollowUseCase {
	return &FollowUseCase{
		followRepo:  followRepo,
		userUseCase: userUseCase,
		tagUseCase:  tagUseCase,
		logger:      logger,
	}
}

func (u *FollowUseCase) Follow(ctx context.Context, follow *followentity.Follow) error {
	if follow.FollowerID == follow.FollowingID {
		return errorx.BadRequest(reason.DisallowFollow)
	}
	return u.followRepo.Follow(ctx, follow)
}

func (u *FollowUseCase) FollowCancel(ctx context.Context, follow *followentity.Follow) error {
	return u.followRepo.FollowCancel(ctx, follow)
}

func (u *FollowUseCase) GetFollowers(ctx context.Context, followQuery *followentity.FollowQuery) (paginator.Page[int64], error) {
	return u.followRepo.GetFollowersPage(ctx, followQuery)
}

func (u *FollowUseCase) GetFollowings(ctx context.Context, followQuery *followentity.FollowQuery) (paginator.Page[int64], error) {
	return u.followRepo.GetFollowingsPage(ctx, followQuery)
}

func (u *FollowUseCase) GetFollowingObjects(ctx context.Context, objectIDs []int64, followedType followentity.FollowedType) ([]any, error) {
	switch followedType {
	case followentity.FollowUser:
		users, err := u.userUseCase.GetByUserIDs(ctx, objectIDs)
		if err != nil {
			return nil, err
		}

		objects := make([]any, len(users))
		for i := 0; i < len(users); i++ {
			link, err := u.userUseCase.GetAvatarLink(ctx, users[i].UserID)
			if err != nil {
				return nil, err
			}
			objects[i] = users[i].ToBasicInfo(link)
		}
		return objects, nil

	case followentity.FollowTag:
		tags, err := u.tagUseCase.GetByTagIDs(ctx, objectIDs)
		if err != nil {
			return nil, err
		}

		objects := make([]any, len(tags))
		for i := 0; i < len(tags); i++ {
			objects[i] = tags[i].ToBasicInfo()
		}
		return objects, nil
	}

	return nil, nil
}

func (u *FollowUseCase) IsFollowing(ctx context.Context, follow *followentity.Follow) (bool, error) {
	return u.followRepo.IsFollowing(ctx, follow)
}

func (u *FollowUseCase) AreFollowEachOther(ctx context.Context, userIDx int64, userIDy int64) (bool, error) {
	return u.followRepo.AreFollowEachOther(ctx, userIDx, userIDy)
}
