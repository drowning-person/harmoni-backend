package main

import (
	"harmoni/internal/conf"
	"harmoni/internal/cron"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	*fiber.App
}

func NewApplication(fiberApp *fiber.App, cronApp *cron.ScheduledTaskManager) *Application {
	cronApp.Start()
	return &Application{App: fiberApp}
}

func main() {
	cfg, err := conf.ReadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	app, clean, err := initApplication(cfg.App, cfg.DB, cfg.Redis, cfg.Auth, cfg.Email, cfg.MessageQueue, cfg.FileStorage, cfg.Log)
	defer clean()
	if err != nil {
		panic(err)
	}

	err = app.Listen(cfg.App.Addr)
	if err != nil {
		panic(err)
	}
}
