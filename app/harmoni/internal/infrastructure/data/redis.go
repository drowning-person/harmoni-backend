package data

import (
	"context"
	"fmt"
	"harmoni/app/harmoni/internal/infrastructure/config"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var CacheProvider = wire.NewSet(
	NewRedis,
	wire.Bind(new(redis.UniversalClient), new(*redis.Client)),
)

func NewRedis(conf *config.Redis) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		Password:     conf.Password,
		DB:           int(conf.Database),
		PoolSize:     conf.PoolSize,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	})

	cleanFunc := func() {
		rdb.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), rdb.Options().ReadTimeout)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, cleanFunc, err
	}

	return rdb, cleanFunc, nil
}
