package main

import (
	"fiberLearn/pkg/snowflake"
	"fiberLearn/router"
)

func main() {
	app := router.New()
	if err := snowflake.Init("2022-07-31", 1); err != nil {
		panic(err)
	}
	app.Listen(":3000")
}
