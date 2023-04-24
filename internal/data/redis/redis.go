package redis

import (
	"context"
	"fmt"
	"harmoni/internal/conf"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
)

func NewRedis(conf *conf.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		Password:     conf.Password,
		DB:           int(conf.Database),
		PoolSize:     conf.PoolSize,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
