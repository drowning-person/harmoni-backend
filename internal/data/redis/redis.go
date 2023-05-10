package redis

import (
	"context"
	"fmt"
	"harmoni/internal/conf"

	"github.com/redis/go-redis/v9"
)

func NewRedis(conf *conf.Redis) (*redis.Client, error) {
	fmt.Printf("redis conf:%#v\n", conf)

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		Password:     conf.Password,
		DB:           int(conf.Database),
		PoolSize:     conf.PoolSize,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), rdb.Options().ReadTimeout)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
