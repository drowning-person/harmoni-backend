//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"harmoni/app/harmoni/internal/cron"
	"harmoni/app/harmoni/internal/handler"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/infrastructure/data"
	"harmoni/app/harmoni/internal/infrastructure/mq/publisher"
	"harmoni/app/harmoni/internal/pkg/app"
	"harmoni/app/harmoni/internal/pkg/logger"
	"harmoni/app/harmoni/internal/pkg/middleware"
	"harmoni/app/harmoni/internal/pkg/snowflakex"
	"harmoni/app/harmoni/internal/repository"
	"harmoni/app/harmoni/internal/server/http"
	"harmoni/app/harmoni/internal/server/mq"
	"harmoni/app/harmoni/internal/service"
	"harmoni/app/harmoni/internal/usecase"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func sugar(l *zap.Logger) *zap.SugaredLogger {
	return l.Sugar()
}

func initApplication(
	conf *config.Config,
	appConf *config.App,
	dbconf *config.DB,
	rdbconf *config.Redis,
	authConf *config.Auth,
	emailConf *config.Email,
	likeConf *config.Like,
	messageConf *config.MessageQueue,
	fileConf *config.FileStorage,
	logConf *config.Log) (*app.Application, func(), error) {
	panic(wire.Build(
		sugar,
		middleware.NewJwtAuthMiddleware,
		http.ProviderSetHTTP,
		snowflakex.NewSnowflakeNode,
		logger.NewZapLogger,
		repository.ProviderSetRepo,
		handler.ProviderSetHandler,
		service.ProviderSetService,
		usecase.ProviderSetUsecase,
		cron.NewScheduledTaskManager,
		newApplication,
		data.NewDB,
		data.CacheProvider,
		mq.ProviderSetMQ,
		publisher.ProviderSetPublisher,
	))
}
