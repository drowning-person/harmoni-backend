package like

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	objectv1 "harmoni/api/common/object/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	polike "harmoni/app/like/internal/infrastructure/po/like"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/consts"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func byObject(object *objectv1.Object) data.ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("object_id = ?", object.GetId()).
			Where("object_type = ?", object.GetType())
	}
}

func genLikeCountKey(objectID int64, objectType objectv1.ObjectType) string {
	return fmt.Sprintf("like:%s:%d", objectType.Format(), objectID)
}

const (
	likeCountKeyTTL = 3600 * time.Second
)

func genLikeCountKeyTTL() time.Duration {
	return likeCountKeyTTL
}

func (r *LikeRepo) getObjectLikeCountFromDB(ctx context.Context, object *objectv1.Object) (int64, error) {
	count := 0
	err := r.data.DB(ctx).Model(&polike.LikeCount{}).
		Select("counts").
		Scopes(byObject(object)).
		First(&count).Error
	if err != nil {
		return 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return int64(count), nil
}

func (r *LikeRepo) getObjectLikeCountFromCache(ctx context.Context, object *objectv1.Object) (int64, error) {
	key := genLikeCountKey(object.GetId(), object.GetType())
	countStr, err := r.rdb.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		count, err := r.getObjectLikeCountFromDB(ctx, object)
		if err != nil {
			return 0, err
		}
		err = r.rdb.Set(ctx, key, count, genLikeCountKeyTTL()).Err()
		if err != nil {
			return 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return count, nil
	case err != nil:
		return 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		count, err := strconv.ParseInt(countStr, consts.BaseDecimal, consts.BitSize64)
		if err != nil {
			return 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return count, nil
	}
}

func (r *LikeRepo) ObjectLikeCount(ctx context.Context, object *objectv1.Object) (*entitylike.LikeCount, error) {
	count, err := r.getObjectLikeCountFromCache(ctx, object)
	if err != nil {
		return nil, err
	}
	return &entitylike.LikeCount{
		Count:  count,
		Object: object,
	}, nil
}

func (r *LikeRepo) ListObjectLikeCount(ctx context.Context, objectIDs []int64, objectType objectv1.ObjectType) (entitylike.LikeCountList, error) {
	lcList := make([]*polike.LikeCount, 0, 10)
	err := r.data.DB(ctx).Model(lcList).
		Where("object_id IN ?", objectIDs).
		Where("object_type = ?", objectType).
		Find(lcList).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return lo.Map(lcList, func(lc *polike.LikeCount, _ int) *entitylike.LikeCount {
		return &entitylike.LikeCount{
			Count: lc.Counts,
			Object: &objectv1.Object{
				Id:   lc.ObjectID,
				Type: lc.OjbectType,
			},
		}
	}), nil
}

func (r *LikeRepo) AddLikeCount(ctx context.Context, object *objectv1.Object, count int64) error {
	err := r.data.DB(ctx).Model(&polike.LikeCount{}).
		Scopes(byObject(object)).
		UpdateColumn("counts", gorm.Expr("counts + ?", count)).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
