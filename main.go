package main

import (
	"fiberLearn/model"
	"fiberLearn/model/redis"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/pkg/validator"
	"fiberLearn/pkg/zap"
	"fiberLearn/router"
)

func main() {
	if err := zap.InitLogger("./log/main.log"); err != nil {
		panic(err)
	}
	if err := model.InitMysql("utf8mb4"); err != nil {
		panic(err)
	}
	if err := snowflake.Init("2022-07-31", 1); err != nil {
		panic(err)
	}
	if err := validator.InitTrans("zh"); err != nil {
		panic(err)
	}
	app := router.New()
	redis.InitRedis()
	app.Listen(":3000")
}
