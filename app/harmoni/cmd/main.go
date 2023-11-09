package main

import (
	"context"
	"encoding/json"
	"harmoni/app/harmoni/internal/cron"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/pkg/app"
	"harmoni/app/harmoni/internal/pkg/validator"
	"harmoni/app/harmoni/internal/server/http"
	"harmoni/app/harmoni/internal/server/mq"

	"go.uber.org/zap"
)

func newApplication(
	conf *config.Config,
	fiberExecutor *http.FiberExecutor,
	cronApp *cron.ScheduledTaskManager,
	messageExecutor *mq.MQExecutor,
	logger *zap.SugaredLogger,
) *app.Application {
	data, err := json.Marshal(conf)
	if err != nil {
		logger.Warnf("print config failed: %s", data)
	} else {
		logger.Infof("config json: %+v", string(data))
	}
	return app.NewApp(logger,
		app.WithServer(fiberExecutor, cronApp, messageExecutor))
}

func main() {
	cfg, err := config.ReadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	app, clean, err := initApplication(cfg, cfg.App, cfg.DB, cfg.Redis, cfg.Auth, cfg.Email, cfg.Like, cfg.MessageQueue, cfg.FileStorage, cfg.Log)
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
