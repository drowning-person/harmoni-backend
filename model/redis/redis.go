package redis

import (
	"context"
	"fmt"
	"harmoni/config"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func InitRedis(conf *config.Redis) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		Password:     conf.Password,
		DB:           int(conf.Database),
		PoolSize:     conf.PoolSize,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
