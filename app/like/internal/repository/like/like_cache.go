package like

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "harmoni/api/common/object/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/internal/pkg/converter"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/consts"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

const (
	cacheUserLikeListCount = 600
)

type likeCacheInfo struct {
	ObjectID  int64
	UpdatedAt time.Time
}

func (l *likeCacheInfo) marshalJSONStr() string {
	data, _ := json.Marshal(l)
	return converter.BytesToString(data)
}

func genUserLikeListKey(objectType v1.ObjectType) string {
	return fmt.Sprintf("user:like:%s", objectType.Format())
}

func (r *LikeRepo) listNLikeObjectByUserIDFromDB(ctx context.Context, n int, query *entitylike.ListLikeObjectQuery) ([]*likeCacheInfo, error) {
	likes := []*likeCacheInfo{}
	err := r.data.DB(ctx).
		Table("like").
		Select("object_id", "updated_at").
		Scopes(withObjectType(query.ObjectType), withUserID(query.UserID)).
		Order("updated_at DESC").
		Limit(n).
		Find(&likes).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return likes, nil
}

func (r *LikeRepo) listLikeObjectByUserIDFromCache(ctx context.Context, query *entitylike.ListLikeObjectQuery) ([]int64, int64, error) {
	key := genUserLikeListKey(query.ObjectType)
	n, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if n == 0 {
		likeList, err := r.listNLikeObjectByUserIDFromDB(ctx, cacheUserLikeListCount, query)
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
		return lo.Map(likeList,
				func(cacheInfo *likeCacheInfo, _ int) int64 { return cacheInfo.ObjectID }),
			int64(len(likeList)), nil
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

func (r *LikeRepo) saveLikeCountToCache(
	ctx context.Context,
	object *v1.Object,
	isCancel bool,
	count int64,
) error {
	key := genLikeCountKey(object)
	if isCancel {
		err := r.rdb.DecrBy(ctx, key, count).Err()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}
	incrValue, err := r.rdb.IncrBy(ctx, key, count).Result()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	// when incrValue == 1, it means the first like
	if incrValue != count {
		return nil
	}
	err = r.rdb.Expire(ctx, key, genLikeCountKeyTTL()).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

const (
	unlikeScript = `
    local removed = redis.call('ZREM', KEYS[1], ARGV[1])
    if removed == 1 then
        return redis.call('DECR', KEYS[2])
    end
	`

	likeScript = `
	local added = redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2])
    if added == 1 then
    	return redis.call('INCR', KEYS[2])
    end
	`
)

func (r *LikeRepo) saveToCache(
	ctx context.Context,
	object *v1.Object,
	isCancel bool,
) error {
	key := genLikeCountKey(object)
	userLikeListKey := genUserLikeListKey(object.Type)
	if isCancel {
		err := redis.NewScript(unlikeScript).Run(ctx, r.rdb,
			[]string{userLikeListKey, key}, object.GetId()).Err()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return nil
	}
	cacheInfo := likeCacheInfo{
		UpdatedAt: time.Now(),
		ObjectID:  object.GetId(),
	}
	incrValue, err := redis.NewScript(likeScript).Run(ctx, r.rdb,
		[]string{userLikeListKey, key},
		float64(cacheInfo.UpdatedAt.Unix()), object.GetId()).Int()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	// when incrValue == 1, it means the first like
	if int64(incrValue) == 1 {
		err = r.rdb.Expire(ctx, key, genLikeCountKeyTTL()).Err()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}
	return nil
}
