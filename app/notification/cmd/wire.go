//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"harmoni/app/notification/internal/handler/remind"
	"harmoni/app/notification/internal/infrastructure/conf"
	"harmoni/app/notification/internal/infrastructure/data"
	"harmoni/app/notification/internal/infrastructure/rpc"
	"harmoni/app/notification/internal/pkg/middleware"
	"harmoni/app/notification/internal/repository"
	"harmoni/app/notification/internal/server/http"
	"harmoni/app/notification/internal/server/mq"
	"harmoni/app/notification/internal/service"
	"harmoni/app/notification/internal/usecase"
	commonconf "harmoni/internal/conf"
	commondata "harmoni/internal/pkg/data"
	"harmoni/internal/pkg/etcdx"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/server"
	"harmoni/internal/pkg/snowflakex"
	commonrepo "harmoni/internal/repository"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

func initApplication(
	conf *conf.Conf,
	appConf *commonconf.App,
	dataConf *commonconf.DB,
	etcdConf *commonconf.ETCD,
	serverConf *commonconf.Server,
	logConf *commonconf.Log,
	mqConf *commonconf.MessageQueue,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		logger.ProviderSetLogger,
		data.NewDB,
		commondata.ProviderSet,
		commonrepo.ProviderSet,
		snowflakex.NewSnowflakeNode,
		etcdx.NewETCDClient,
		repository.ProviderSet,
		usecase.ProviderSet,
		service.ProviderSetService,
		rpc.NewUsergRPC,
		middleware.NewAuthUserMiddleware,
		mq.NewMQRouter,
		remind.NewRemindHandler,
		http.ProviderSetHTTP,
		server.NewMQServer,
		newApplication,
	))
}
