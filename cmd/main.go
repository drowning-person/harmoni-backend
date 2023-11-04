package main

import (
	"context"
	"harmoni/internal/cron"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/pkg/app"
	"harmoni/internal/pkg/validator"
	"harmoni/internal/server/http"
	"harmoni/internal/server/mq"

	"go.uber.org/zap"
)

func newApplication(
	fiberExecutor *http.FiberExecutor,
	cronApp *cron.ScheduledTaskManager,
	messageExecutor *mq.MQExecutor,
	logger *zap.SugaredLogger,
) *app.Application {
	return app.NewApp(logger,
		app.WithServer(fiberExecutor, cronApp, messageExecutor))
}

func main() {
	cfg, err := config.ReadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	app, clean, err := initApplication(cfg.App, cfg.DB, cfg.Redis, cfg.Auth, cfg.Email, cfg.MessageQueue, cfg.FileStorage, cfg.Log)
	defer clean()
	if err != nil {
		panic(err)
	}
	err = validator.InitTrans(cfg.App)
	if err != nil {
		panic(err)
	}

	err = app.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
