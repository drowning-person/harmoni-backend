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
	"harmoni/app/harmoni/internal/pkg/middleware"
	"harmoni/app/harmoni/internal/pkg/snowflakex"
	"harmoni/app/harmoni/internal/repository"
	"harmoni/app/harmoni/internal/server/grpc"
	"harmoni/app/harmoni/internal/server/http"
	"harmoni/app/harmoni/internal/server/mq"
	"harmoni/app/harmoni/internal/service"
	"harmoni/app/harmoni/internal/usecase"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/etcdx"
	"harmoni/internal/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func sugar(l *zap.Logger) *zap.SugaredLogger {
	return l.Sugar()
}

func initApplication(
	conf *config.Config,
	appConf *config.App,
	dbconf *conf.Database,
	rdbconf *config.Redis,
	authConf *config.Auth,
	emailConf *config.Email,
	likeConf *config.Like,
	messageConf *config.MessageQueue,
	fileConf *config.FileStorage,
	etcdConf *conf.ETCD,
	serverConf *conf.Server,
	logConf *conf.Log) (*kratos.App, func(), error) {
	panic(wire.Build(
		sugar,
		logger.ProviderSetLogger,
		snowflakex.NewSnowflakeNode,
		data.NewDB,
		data.CacheProvider,
		etcdx.NewETCDClient,
		middleware.NewJwtAuthMiddleware,
		http.ProviderSetHTTP,
		grpc.NewGrpcServer,
		mq.ProviderSetMQ,
		repository.ProviderSetRepo,
		handler.ProviderSetHandler,
		service.ProviderSetService,
		usecase.ProviderSetUsecase,
		cron.NewScheduledTaskManager,
		publisher.ProviderSetPublisher,
		newApplication,
	))
}
