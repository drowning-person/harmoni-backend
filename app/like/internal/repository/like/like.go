package like

import (
	"context"
	"errors"

	entitylike "harmoni/app/like/internal/entity/like"
	polike "harmoni/app/like/internal/infrastructure/po/like"
	reasonlike "harmoni/app/like/internal/pkg/reason"
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
			err := tx.Scopes(findLike(like)).
				Delete(&polike.Like{}).Error
			if err != nil {
				return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			return nil
		}

		po := polike.FromDomain(like)
		model := &polike.Like{}
		err := tx.Model(model).Unscoped().Scopes(findLike(like)).First(model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Table(model.TableName()).
				Update("deleted_at", 0).
				Where("liking_id = ?", po.LikingID).Error
			if err != nil {
				return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			return nil
		} else if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

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
	err := r.data.DB(ctx).Model(&polike.Like{}).
		Scopes(findLike(like)).
		Count(&count).Error
	if err != nil {
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count > 0, nil
}

func (r *LikeRepo) Get(ctx context.Context, like *entitylike.Like) (*entitylike.Like, error) {
	l := &entitylike.Like{}
	err := r.data.DB(ctx).Model(&polike.Like{}).
		Scopes(findLike(like)).
		First(l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound(reasonlike.LikeNotExist)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return l, nil
}
