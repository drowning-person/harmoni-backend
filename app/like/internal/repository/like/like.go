package like

import (
	"context"

	entitylike "harmoni/app/like/internal/entity/like"
	polike "harmoni/app/like/internal/infrastructure/po/like"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/iface"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type LikeRepo struct {
	uniqueRepo iface.UniqueIDRepository
	data       *data.DB
	logger     *log.Helper
}

var _ entitylike.LikeRepository = (*LikeRepo)(nil)

func findLike(like *entitylike.Like) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", like.User.GetId()).
			Where("target_user_id = ?", like.TargetUser.GetId()).
			Where("like_type = ?", like.LikeType).
			Where("object_id = ?", like.ObjectID)
	}
}

func NewLikeRepo(
	uniqueRepo iface.UniqueIDRepository,
	data *data.DB,
	logger log.Logger,
) *LikeRepo {
	return &LikeRepo{
		uniqueRepo: uniqueRepo,
		data:       data,
		logger: log.NewHelper(
			log.With(logger, "module", "repository/like", "service", "like")),
	}
}

func (r *LikeRepo) Save(ctx context.Context, like *entitylike.Like, isCancel bool) error {
	return r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if isCancel {
			err := tx.Table("like").Scopes(findLike(like)).
				Delete(&polike.Like{}).Error
			if err != nil {
				return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			return nil
		}
		po := polike.FromDomain(like)
		var err error
		po.LikingID, err = r.uniqueRepo.GenUniqueID(ctx)
		if err != nil {
			return err
		}
		err = tx.Table("like").Create(po).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return nil
	})
}

func (r *LikeRepo) IsExist(ctx context.Context, like *entitylike.Like) (bool, error) {
	var count int64
	err := r.data.DB(ctx).Table("like").
		Scopes(findLike(like)).
		Count(&count).Error
	if err != nil {
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count > 0, nil
}
