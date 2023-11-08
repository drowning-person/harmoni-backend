//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"harmoni/internal/cron"
	"harmoni/internal/handler"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/infrastructure/data"
	"harmoni/internal/infrastructure/mq/publisher"
	"harmoni/internal/pkg/app"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/repository"
	"harmoni/internal/server/http"
	"harmoni/internal/server/mq"
	"harmoni/internal/service"
	"harmoni/internal/usecase"

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
