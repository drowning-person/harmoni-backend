package follow

import (
	"context"
	followentity "harmoni/app/harmoni/internal/entity/follow"
	"harmoni/app/harmoni/internal/entity/paginator"
	"harmoni/app/harmoni/internal/entity/tag"
	"harmoni/app/harmoni/internal/entity/unique"
	"harmoni/app/harmoni/internal/entity/user"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ followentity.FollowRepository = (*FollowRepo)(nil)

type FollowRepo struct {
	db           *gorm.DB
	uniqueIDRepo unique.UniqueIDRepo
	userRepo     user.UserRepository
	tagRepo      tag.TagRepository
	logger       *zap.SugaredLogger
}

func NewFollowRepo(db *gorm.DB,
	userRepo user.UserRepository,
	tagRepo tag.TagRepository,
	uniqueIDRepo unique.UniqueIDRepo,
	logger *zap.SugaredLogger) *FollowRepo {
	return &FollowRepo{
		db:           db,
		userRepo:     userRepo,
		tagRepo:      tagRepo,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/follow"),
	}
}

func (r *FollowRepo) Follow(ctx context.Context, follow *followentity.Follow) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		exist, err := r.getFollowObject(ctx, tx, follow)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if !exist {
			return errorx.NotFound(reason.ObjectNotFound)
		}

		exist, err = r.isFollowExist(ctx, tx, follow)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if exist {
			return errorx.BadRequest(reason.FollowAlreadyExist)
		}

		err = tx.Create(follow).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		err = r.updateFollows(ctx, tx, follow, 1)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *FollowRepo) FollowCancel(ctx context.Context, follow *followentity.Follow) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		exist, err := r.getFollowObject(ctx, tx, follow)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if !exist {
			return errorx.NotFound(reason.ObjectNotFound)
		}

		exist, err = r.isFollowExist(ctx, tx, follow)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if !exist {
			return errorx.NotFound(reason.FollowNotFound)
		}

		err = r.db.WithContext(ctx).
			Where("follower_id = ? AND following_id = ? AND followed_type = ?", follow.FollowerID, follow.FollowingID, follow.FollowedType).
			Delete(follow).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		err = r.updateFollows(ctx, tx, follow, -1)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *FollowRepo) getFollowObject(ctx context.Context, tx *gorm.DB, follow *followentity.Follow) (bool, error) {
	var (
		exist bool
		err   error
	)
	switch follow.FollowedType {
	case followentity.FollowUser:
		_, exist, err = r.userRepo.GetByUserID(ctx, follow.FollowingID)
	case followentity.FollowTag:
		_, exist, err = r.tagRepo.GetByTagID(ctx, follow.FollowingID)
	}

	return exist, err
}

func (r *FollowRepo) isFollowExist(ctx context.Context, tx *gorm.DB, follow *followentity.Follow) (bool, error) {
	var count int64
	err := tx.Model(follow).
		Where("follower_id = ? AND following_id = ? AND followed_type = ?", follow.FollowerID, follow.FollowingID, follow.FollowedType).
		Count(&count).Error
	return count > 0, err
}

func (r *FollowRepo) updateFollows(ctx context.Context, tx *gorm.DB, follow *followentity.Follow, followCount int) error {
	var err error
	switch follow.FollowedType {
	case followentity.FollowUser:
		err = tx.Table("user").Where("user_id = ?", follow.FollowingID).UpdateColumn("follow_count", gorm.Expr("follow_count + ? ", followCount)).Error
	case followentity.FollowTag:
		err = tx.Table("tag").Where("tag_id = ?", follow.FollowingID).UpdateColumn("follow_count", gorm.Expr("follow_count + ? ", followCount)).Error
	default:
		err = errorx.BadRequest(reason.DisallowFollow).WithMsg("this object can't be followed")
	}

	return err
}

func (r *FollowRepo) GetFollowersPage(ctx context.Context, followQuery *followentity.FollowQuery) (paginator.Page[int64], error) {
	idPage := paginator.Page[int64]{
		CurrentPage: followQuery.Page,
		PageSize:    followQuery.PageSize,
		Data:        make([]int64, 0, 8),
	}
	err := r.db.WithContext(ctx).Table("follow").
		Where("following_id = ? AND followed_type = ? AND deleted_at is NULL", followQuery.UserID, followQuery.Type).
		Count(&idPage.Total).Error
	if err != nil {
		return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	err = r.db.WithContext(ctx).Table("follow").
		Select("follower_id").
		Where("following_id = ? AND followed_type = ? AND deleted_at is NULL", followQuery.UserID, followQuery.Type).
		Scopes(paginator.Paginate(&idPage)).
		Find(&idPage.Data).Error

	if err != nil {
		return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return idPage, nil
}

func (r *FollowRepo) GetFollowingsPage(ctx context.Context, followQuery *followentity.FollowQuery) (paginator.Page[int64], error) {
	idPage := paginator.Page[int64]{
		CurrentPage: followQuery.Page,
		PageSize:    followQuery.PageSize,
		Data:        make([]int64, 0, 8),
	}
	err := r.db.WithContext(ctx).Table("follow").
		Where("follower_id = ? AND followed_type = ? AND deleted_at is NULL", followQuery.UserID, followQuery.Type).
		Count(&idPage.Total).Error
	if err != nil {
		return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	err = paginator.Paginate(&idPage)(r.db.WithContext(ctx)).Table("follow").
		Select("following_id").
		Where("follower_id = ? AND followed_type = ? AND deleted_at is NULL", followQuery.UserID, followQuery.Type).
		Scopes(paginator.Paginate(&idPage)).
		Find(&idPage.Data).Error

	if err != nil {
		return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return idPage, nil
}

func (r *FollowRepo) GetFollowingUsersAll(ctx context.Context, userID int64) ([]int64, error) {
	followings := []int64{}
	err := r.db.WithContext(ctx).Table("follow").
		Select("following_id").
		Where("follower_id = ? AND followed_type = ? AND deleted_at is NULL", userID, followentity.FollowUser).
		Find(&followings).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return followings, nil
}

func (r *FollowRepo) IsFollowing(ctx context.Context, follow *followentity.Follow) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(follow).
		Where("follower_id = ? AND following_id = ? AND followed_type = ?", follow.FollowerID, follow.FollowingID, follow.FollowedType).
		Count(&count).Error; err != nil {
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return count > 0, nil
}

func (r *FollowRepo) AreFollowEachOther(ctx context.Context, userIDx int64, userIDy int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&followentity.Follow{}).
		Joins("INNER JOIN follow AS f2 ON follow.follower_id = f2.following_id AND follow.following_id = f2.follower_id").
		Where("follow.follower_id = ? AND f2.follower_id = ?", userIDx, userIDy).
		Count(&count).Error; err != nil {
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return count > 0, nil
}
