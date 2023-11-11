package like

import (
	"context"
	"encoding/binary"
	"fmt"
	"harmoni/app/harmoni/internal/entity"
	likeentity "harmoni/app/harmoni/internal/entity/like"
	"harmoni/app/harmoni/internal/entity/paginator"
	userentity "harmoni/app/harmoni/internal/entity/user"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ likeentity.LikeRepository = (*LikeRepo)(nil)

const (
	shardCounts = 5
)

const (
	postLikeCountPrefix    = "post:like.count:"
	commentLikeCountPrefix = "comment:like.count:"
	userLikeCountPrefix    = "user:like.count:"
)

func hashID(id int64) uint32 {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(id))
	return crc32.ChecksumIEEE(bs) % shardCounts
}

func postLikeCountKey(id int64) string {
	hashedID := hashID(id)
	return fmt.Sprintf("%s%d", postLikeCountPrefix, hashedID)
}

func commentLikeCountKey(id int64) string {
	hashedID := hashID(id)
	return fmt.Sprintf("%s%d", commentLikeCountPrefix, hashedID)
}

func userLikeCountKey(id int64) string {
	hashedID := hashID(id)
	return fmt.Sprintf("%s%d", userLikeCountPrefix, hashedID)
}

func userLikingSetKey(id int64, likeType likeentity.LikeType) string {
	return fmt.Sprintf("user:%d:like.recent:%d", id, likeType)
}

func getCountCacheKey(like *likeentity.Like) string {
	var key string
	switch like.LikeType {
	case likeentity.LikePost:
		key = postLikeCountKey(like.LikingID)
	case likeentity.LikeComment:
		key = commentLikeCountKey(like.LikingID)
	case likeentity.LikeUser:
		key = userLikeCountKey(like.LikingID)
	}

	return key
}

type LikeRepo struct {
	conf     *config.Like
	db       *gorm.DB
	rdb      redis.UniversalClient
	userRepo userentity.UserRepository
	logger   *zap.SugaredLogger
}

func NewLikeRepo(
	conf *config.Like,
	db *gorm.DB,
	rdb redis.UniversalClient,
	userRepo userentity.UserRepository,
	logger *zap.SugaredLogger,
) *LikeRepo {
	return &LikeRepo{
		conf:     conf,
		db:       db,
		rdb:      rdb,
		userRepo: userRepo,
		logger:   logger.With("module", "repository/like"),
	}
}

func (r *LikeRepo) Save(ctx context.Context, like *likeentity.Like) error {
	err := r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "liking_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"canceled"}),
		},
	).Create(like).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return err
}

/*
	 func (r *LikeRepo) UpdateLikeCount(ctx context.Context, like *likeentity.Like, count int8) error {
		objKey := getCountCacheKey(like)
		_, err := r.rdb.HIncrBy(ctx, objKey, strconv.FormatInt(like.LikingID, 10), int64(count)).Result()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return nil
	}
*/

func (r *LikeRepo) cacheHSetLikeCount(ctx context.Context, key string, id int64, count interface{}) error {
	// set count to cache
	err := r.rdb.HSet(ctx, key, id, count).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	// expire cache
	err = r.rdb.ExpireNX(ctx, key, r.conf.CacheDuration).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

const (
	unlikeScript = `
	local liked = redis.call('ZSCORE', KEYS[1], ARGV[1])
	if not liked then
		return 0
	end
    local removed = redis.call('ZREM', KEYS[1], ARGV[1])
    if removed == 1 then
        redis.call('HINCRBY', KEYS[2], ARGV[1], ARGV[3])
		redis.call('HINCRBY', KEYS[3], ARGV[4], ARGV[3])
    end
	return 1
	`

	likeScript = `
	local liked = redis.call('ZSCORE', KEYS[1], ARGV[1])
	if liked then
		return 0
	end
	local added = redis.call('ZADD', KEYS[1], ARGV[2], ARGV[1])
    if added == 1 then
    	redis.call('HINCRBY', KEYS[2], ARGV[1], ARGV[3])
		redis.call('HINCRBY', KEYS[3], ARGV[4], ARGV[3])
    end
	return 1
	`
)

func (r *LikeRepo) likeAction(ctx context.Context, key string, like *likeentity.Like, targetUserID int64, isCancel bool) error {
	var (
		count int64
	)

	objKey := getCountCacheKey(like)
	if n, err := r.rdb.Exists(ctx, objKey).Result(); err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if n < 1 {
		// retrieve count from DB
		err = r.db.WithContext(ctx).Table("like").Where("liking_id = ? AND like_type = ? AND canceled = 0", like.LikingID, like.LikeType).Count(&count).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		err = r.cacheHSetLikeCount(ctx, objKey, like.LikingID, count)
		if err != nil {
			return err
		}
	}
	userLikedKey := userLikeCountKey(targetUserID)
	err := r.rdb.HGet(ctx, userLikedKey, strconv.FormatInt(targetUserID, 10)).Err()
	if err == redis.Nil {
		// user is impossible not exist
		count, _, err := r.userRepo.GetLikeCount(ctx, targetUserID)
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		// set count to cache
		err = r.cacheHSetLikeCount(ctx, userLikedKey, targetUserID, count)
		if err != nil {
			return err
		}
	} else if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if isCancel {
		score, err := redis.NewScript(unlikeScript).Run(ctx, r.rdb, []string{key, objKey, userLikedKey}, like.LikingID, time.Now().Unix(), -1, targetUserID).Int()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if score == 0 {
			return errorx.BadRequest(reason.LikeCancelFailToNotLiked)
		}
	} else {
		score, err := redis.NewScript(likeScript).Run(ctx, r.rdb, []string{key, objKey, userLikedKey}, like.LikingID, time.Now().Unix(), 1, targetUserID).Int()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if score == 0 {
			return errorx.BadRequest(reason.LikeAlreadyExist)
		}
	}

	return nil
}

func (r *LikeRepo) transLikingsFromDBToCache(ctx context.Context, key string, like *likeentity.Like) ([]likeentity.LikeCacheInfo, error) {
	// user recent likes
	likes := []likeentity.LikeCacheInfo{}
	err := r.db.WithContext(ctx).
		Table("like").
		Select("liking_id", "updated_at").
		Where("user_id = ? AND like_type = ? AND canceled = 0", like.UserID, like.LikeType).
		Order("updated_at DESC").
		Limit(600).
		Find(&likes).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	memebers := make([]redis.Z, len(likes)+1)
	for i, liking := range likes {
		memebers[i] = redis.Z{
			Score:  float64(liking.UpdatedAt),
			Member: liking.LikingID,
		}
	}
	memebers[len(likes)] = redis.Z{Score: -1, Member: entity.DefaultRedisValue}
	err = r.rdb.ZAddArgs(ctx, key, redis.ZAddArgs{
		NX:      true,
		Members: memebers,
	}).Err()
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	expired, err := r.rdb.ExpireNX(ctx, key, r.conf.CacheDuration).Result()
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if !expired {
		r.logger.Warnf("expire key:%s failed", key)
	}

	return likes, nil
}

func (r *LikeRepo) Like(ctx context.Context, like *likeentity.Like, targetUserID int64, isCancel bool) error {
	var (
		err error
	)

	key := userLikingSetKey(like.UserID, like.LikeType)
	if n, err := r.rdb.Exists(ctx, key).Result(); err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if n > 0 {
		err = r.likeAction(ctx, key, like, targetUserID, isCancel)
		if err != nil {
			return err
		}
		return nil
	}

	// user recent likes
	_, err = r.transLikingsFromDBToCache(ctx, key, like)
	if err != nil {
		return err
	}

	err = r.likeAction(ctx, key, like, targetUserID, isCancel)
	if err != nil {
		return err
	}

	return nil
}

func (r *LikeRepo) CacheLikeCount(ctx context.Context, like *likeentity.Like, count int64) error {
	key := getCountCacheKey(like)

	var value interface{}
	if count == -1 {
		value = entity.DefaultRedisValue
	} else {
		value = count
	}

	return r.cacheHSetLikeCount(ctx, key, like.LikingID, value)
}

func (r *LikeRepo) LikeCount(ctx context.Context, like *likeentity.Like) (int64, bool, error) {
	key := getCountCacheKey(like)

	countStr, err := r.rdb.HGet(ctx, key, strconv.FormatInt(like.LikingID, 10)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	var count int64
	if countStr != entity.DefaultRedisValue {
		count, err = strconv.ParseInt(countStr, 10, 64)
	} else {
		return 0, false, errorx.NotFound(reason.ObjectNotFound)
	}

	return count, true, err
}

func (r *LikeRepo) BatchLikeCount(ctx context.Context, likeType likeentity.LikeType) (map[int64]int64, error) {
	keys := make([]string, 0, shardCounts)
	switch likeType {
	case likeentity.LikePost:
		for i := 0; i < shardCounts; i++ {
			keys = append(keys, fmt.Sprintf("%s%d", postLikeCountPrefix, i))
		}
	case likeentity.LikeComment:
		for i := 0; i < shardCounts; i++ {
			keys = append(keys, fmt.Sprintf("%s%d", commentLikeCountPrefix, i))
		}
	case likeentity.LikeUser:
		for i := 0; i < shardCounts; i++ {
			keys = append(keys, fmt.Sprintf("%s%d", userLikeCountPrefix, i))
		}
	}

	m := map[int64]int64{}
	for _, key := range keys {
		kvs, err := r.rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		for k, v := range kvs {
			if v == entity.DefaultRedisValue {
				continue
			}
			vNum, _ := strconv.ParseInt(v, 10, 64)

			kNum, _ := strconv.ParseInt(k, 10, 64)
			m[kNum] = vNum
		}
	}

	return m, nil
}

func (r *LikeRepo) BatchLikeCountByIDs(ctx context.Context, likingIDs []int64, likeType likeentity.LikeType) (map[int64]int64, error) {
	keys := make(map[string][]string)
	countMap := make(map[int64]int64)

	switch likeType {
	case likeentity.LikePost:
		for _, likingID := range likingIDs {
			hashedID := postLikeCountKey(likingID)
			keys[hashedID] = append(keys[hashedID], strconv.FormatInt(likingID, 10))
		}
	case likeentity.LikeComment:
		for _, likingID := range likingIDs {
			hashedID := commentLikeCountKey(likingID)
			keys[hashedID] = append(keys[hashedID], strconv.FormatInt(likingID, 10))
		}
	}

	notCachedIDMap := make(map[string][]string)
	for key, ids := range keys {
		counts, err := r.rdb.HMGet(ctx, key, ids...).Result()
		if err != nil {
			return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		for i, idStr := range ids {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			value := counts[i]
			if value == nil {
				notCachedIDMap[key] = append(notCachedIDMap[key], idStr)
			}
			switch t := value.(type) {
			case string:
				if t != entity.DefaultRedisValue {
					likeCount, _ := strconv.ParseInt(t, 10, 64)
					countMap[id] = likeCount
				}
			}
		}
	}

	var likeCounts []struct {
		LikingID  int64 `json:"liking_id,omitempty"`
		LikeCount int64 `json:"like_count,omitempty"`
	}

	if len(notCachedIDMap) == 0 {
		return countMap, nil
	}

	notCacheIDs := make([]string, 0, 24)
	for _, v := range notCachedIDMap {
		notCacheIDs = append(notCacheIDs, v...)
	}
	err := r.db.WithContext(ctx).Table("like").
		Select("liking_id", "COUNT(*) AS like_count").
		Where("liking_id IN (?) AND canceled = 0", notCacheIDs).
		Group("liking_id").Find(&likeCounts).Error

	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	keyvals := make(map[string][]interface{})
	likeCountsMap := make(map[string]int64)
	for _, v := range likeCounts {
		countMap[v.LikingID] = v.LikeCount
		likeCountsMap[strconv.FormatInt(v.LikingID, 10)] = v.LikeCount
	}
	for key, ids := range notCachedIDMap {
		keyvals[key] = make([]interface{}, 0, len(ids)*2)
		for _, id := range ids {
			keyvals[key] = append(keyvals[key], id, likeCountsMap[id])
		}
		err = r.rdb.HMSet(ctx, key, keyvals[key]...).Err()
		if err != nil {
			return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		err = r.rdb.Expire(ctx, key, r.conf.CacheDuration).Err()
		if err != nil {
			return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}

	return countMap, nil
}

func (r *LikeRepo) paginate(page *entity.PageCond) (int64, int64) {
	if page.Page < 1 || page.Page > 30 {
		page.Page = 1
	}
	if page.PageSize < 1 || page.PageSize > 25 {
		page.PageSize = 25
	}
	start := (page.Page - 1) * page.PageSize
	return start, start + page.PageSize
}

// ListLikingIDs list recent 600 likes of user
func (r *LikeRepo) ListLikingIDs(ctx context.Context, query *likeentity.LikeQuery) (paginator.Page[int64], error) {
	idPage := paginator.Page[int64]{
		CurrentPage: query.Page,
		PageSize:    query.PageSize,
	}
	start, end := r.paginate(&query.PageCond)
	idPage.PageSize = end - start
	if idPage.CurrentPage <= 0 {
		idPage.CurrentPage = 1
	}

	key := userLikingSetKey(query.UserID, query.Type)
	if n, err := r.rdb.Exists(ctx, key).Result(); err != nil {
		return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if n > 0 {
		ids, err := r.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:   key,
			Start: start,
			Stop:  end - 1,
			Rev:   true,
		}).Result()
		if err != nil {
			return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		total, err := r.rdb.ZCard(ctx, key).Result()
		if err != nil {
			return paginator.Page[int64]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		idPage.Data = make([]int64, 0, len(ids))
		for _, idStr := range ids {
			if idStr == entity.DefaultRedisValue {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err == nil {
				idPage.Data = append(idPage.Data, id)
			} else {
				r.logger.Warn(err)
			}
		}
		idPage.Total = total - 1

		idPage.Pages = idPage.Total / idPage.PageSize
		if idPage.Total%idPage.PageSize != 0 {
			idPage.Pages++
		}
		if idPage.CurrentPage > idPage.Pages {
			idPage.CurrentPage = idPage.Pages
		}
		return idPage, nil
	}

	likes, err := r.transLikingsFromDBToCache(ctx, key, &likeentity.Like{UserID: query.UserID, LikingID: query.ObjectID, LikeType: query.Type})
	if err != nil {
		return paginator.Page[int64]{}, err
	}

	idPage.Total = int64(len(likes))
	idPage.Pages = idPage.Total / idPage.PageSize
	if idPage.Total%idPage.PageSize != 0 {
		idPage.Pages++
	}
	if idPage.CurrentPage > idPage.Pages {
		idPage.CurrentPage = idPage.Pages
	}
	if idPage.Total > 0 {
		offset := (idPage.CurrentPage - 1) * idPage.PageSize
		var tmp []likeentity.LikeCacheInfo
		if offset+idPage.PageSize > idPage.Total {
			tmp = likes[offset:idPage.Total]
		} else {
			tmp = likes[offset : offset+idPage.PageSize]
		}
		idPage.Data = make([]int64, len(tmp))
		for i, v := range tmp {
			idPage.Data[i] = v.LikingID
		}
	}

	return idPage, nil
}

func (r *LikeRepo) IsLiking(ctx context.Context, like *likeentity.Like) (bool, error) {
	key := userLikingSetKey(like.UserID, like.LikeType)
	err := r.rdb.ZScore(ctx, key, strconv.FormatInt(like.LikingID, 10)).Err()
	if err == redis.Nil {
		_, err = r.transLikingsFromDBToCache(ctx, key, like)
		if err != nil {
			return false, err
		}
		err = r.rdb.ZScore(ctx, key, strconv.FormatInt(like.LikingID, 10)).Err()
		if err != nil {
			if err == redis.Nil {
				return false, nil
			}
			return false, err
		}
	}
	if err != nil {
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return true, nil
}
