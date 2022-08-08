package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func InitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Woaini.12", // no password set
		DB:       0,           // use default DB
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
