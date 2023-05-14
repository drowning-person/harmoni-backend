//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"harmoni/internal/conf"
	"harmoni/internal/handler"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/repository"
	"harmoni/internal/router"
	"harmoni/internal/server"
	"harmoni/internal/service"
	"harmoni/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func sugar(l *zap.Logger) *zap.SugaredLogger {
	return l.Sugar()
}

func initApplication(appConf *conf.App,
	dbconf *conf.DB,
	rdbconf *conf.Redis,
	authConf *conf.Auth,
	emailConf *conf.Email,
	logConf *conf.Log) (*fiber.App, func(), error) {
	panic(wire.Build(
		// validator.InitTrans("zh"),
		sugar,
		middleware.NewJwtAuthMiddleware,
		router.NewHarmoniAPIRouter,
		snowflakex.NewSnowflakeNode,
		logger.NewZapLogger,
		repository.ProviderSetRepo,
		handler.ProviderSetHandler,
		service.ProviderSetService,
		usecase.ProviderSetUsecase,
		server.NewHTTPServer,
	))
}
