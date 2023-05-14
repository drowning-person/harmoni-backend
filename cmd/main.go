package main

import (
	"harmoni/internal/conf"
)

func main() {
	cfg, err := conf.ReadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	app, clean, err := initApplication(cfg.App, cfg.DB, cfg.Redis, cfg.Auth, cfg.Email, cfg.Log)
	defer clean()
	if err != nil {
		panic(err)
	}

	err = app.Listen(cfg.App.Addr)
	if err != nil {
		panic(err)
	}
}
