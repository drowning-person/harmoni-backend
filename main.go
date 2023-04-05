package main

import (
	"harmoni/config"
	"harmoni/model"
	"harmoni/model/redis"
	"harmoni/pkg/snowflake"
	"harmoni/pkg/validator"
	"harmoni/pkg/zap"
	"harmoni/router"
)

func main() {
	cfg, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	if err := zap.InitLogger(cfg.Log); err != nil {
		panic(err)
	}
	if err := model.InitMysql(cfg.DB); err != nil {
		panic(err)
	}
	if err := snowflake.Init("2022-07-31", 1); err != nil {
		panic(err)
	}
	if err := validator.InitTrans("zh"); err != nil {
		panic(err)
	}
	if err := redis.InitRedis(cfg.Redis); err != nil {
		panic(err)
	}
	app := router.New()

	app.Listen(cfg.App.Addr)
}
