package main

import (
	"harmoni/internal/conf"
)

func main() {
	cfg, err := conf.ReadConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	app, err := initApp(cfg.App, cfg.DB, cfg.Redis, cfg.Auth, cfg.Log)
	if err != nil {
		panic(err)
	}

	err = app.Listen(cfg.App.Addr)
	if err != nil {
		panic(err)
	}
}
