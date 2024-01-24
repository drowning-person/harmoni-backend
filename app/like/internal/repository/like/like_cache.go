package like

import (
	"context"
	"fmt"
	v1 "harmoni/api/common/object/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	polike "harmoni/app/like/internal/infrastructure/po/like"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/consts"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func genUserLikeListKey(objectType v1.ObjectType) string {
	return fmt.Sprintf("user:like:%s", objectType.Format())
}

func (r *LikeRepo) listLikeObjectByUserIDFromDB(ctx context.Context, query *entitylike.ListLikeObjectQuery) ([]*polike.Like, int64, error) {
	likeList := make([]*polike.Like, 0, query.Size)
	err := r.data.DB(ctx).Scopes(
		withUserID(query.UserID),
		withObjectType(query.ObjectType),
	).Find(&likeList).Error
	if err != nil {
		r.logger.Error(err)
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	var count int64
	err = r.data.DB(ctx).Model(&polike.Like{}).Scopes(
		withUserID(query.UserID),
		withObjectType(query.ObjectType),
	).Count(&count).Error
	if err != nil {
		r.logger.Error(err)
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return likeList, count, nil
}

func (r *LikeRepo) listLikeObjectByUserIDFromCache(ctx context.Context, query *entitylike.ListLikeObjectQuery) ([]int64, int64, error) {
	key := genUserLikeListKey(query.ObjectType)
	n, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if n == 0 {
		likeList, total, err := r.listLikeObjectByUserIDFromDB(ctx, query)
		if err != nil {
			return nil, 0, err
		}
		members := make([]redis.Z, 0, len(likeList))
		for _, like := range likeList {
			members = append(members, redis.Z{
				Score:  float64(like.UpdatedAt.Unix()),
				Member: like.ObjectID,
			})
		}
		err = r.rdb.ZAdd(ctx, key, members...).Err()
		if err != nil {
			return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return
	}

	idStrs, err := r.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: query.Start(),
		Stop:  query.End(),
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	total, err := r.rdb.ZCard(ctx, key).Result()
	if err != nil {
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, consts.BaseDecimal, consts.BitSize64)
		if err != nil {
			return nil, 0, errorx.InternalServer(reason.ServerError).WithError(err).WithStack()
		}
		ids = append(ids, id)
	}

	return ids, total, nil
}
