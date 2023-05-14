package post

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	// KeyPostTimeZSet key是 帖子id value 是帖子的发布时间
	KeyPostTimeZSet = "bloom:post:time"
	// KeyLikeNumberZSet key是 帖子id value 是帖子的点赞数量
	KeyLikeNumberZSet = "bloom:post:like:number"
	// KeyPostLikeZetPrefix  key是userid value是点赞或者点踩 后面需要拼接帖子的id
	KeyPostLikeZetPrefix = "bloom:post:like:postID:"
)

func getRedisKeyForLikeUserSet(postID int64) string {
	return fmt.Sprintf("%s%d", KeyPostLikeZetPrefix, postID)
}

func (r *PostRepo) getPostLikeNumber(ctx context.Context, postID int64) (int64, error) {
	count, err := r.rdb.ZScore(ctx, KeyLikeNumberZSet, strconv.FormatInt(postID, 10)).Result()
	if err != nil {
		return 0, err
	}

	return int64(count), err
}

// 按照点赞数 降序排列
func (r *PostRepo) getPostIDsByScore(ctx context.Context, pageSize int64, pageNum int64) (ids []string, err error) {
	start := (pageNum - 1) * pageSize
	ids, err = r.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   KeyLikeNumberZSet,
		Start: start,
		Stop:  start + pageSize - 1,
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, err
	}

	return ids, err
}

// AddPost 每次发表帖子成功 都去 zset里面 新增一条记录
func (r *PostRepo) addPost(ctx context.Context, postID int64) error {
	_, err := r.rdb.ZAdd(ctx, KeyLikeNumberZSet, redis.Z{
		Score:  0,
		Member: strconv.FormatInt(postID, 10),
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

// CheckLike 判断之前有没有投过票 true 代表之前 投过 false 代表之前没有投过
func (r *PostRepo) checkLike(ctx context.Context, postID int64, userID int64) (float64, bool, error) {
	like := r.rdb.ZScore(ctx, getRedisKeyForLikeUserSet(postID), strconv.FormatInt(userID, 10))

	res, err := like.Result()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}

	return res, true, nil
}

// DoLike 点赞 或者点踩 记录这个用户对这个帖子的行为
func (r *PostRepo) doLike(ctx context.Context, postID int64, userID int64, direction int8) error {
	if direction == 2 {
		direction = -1
	}

	pipeLine := r.rdb.TxPipeline()
	value := redis.Z{
		Score:  float64(direction),
		Member: strconv.FormatInt(userID, 10),
	}

	pipeLine.ZAdd(ctx, getRedisKeyForLikeUserSet(postID), value)
	pipeLine.ZIncrBy(ctx, KeyLikeNumberZSet, float64(direction), strconv.FormatInt(postID, 10))
	_, err := pipeLine.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
