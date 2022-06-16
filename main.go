package main

import (
	"fiberLearn/router"
)

func main() {
	app := router.New()
	app.Listen(":80")
}
