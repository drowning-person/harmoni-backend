//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"harmoni/app/like/internal/conf"
	"harmoni/app/like/internal/repository"
	"harmoni/app/like/internal/server"
	"harmoni/app/like/internal/service"
	"harmoni/app/like/internal/usecase"
	commonconf "harmoni/internal/conf"
	commondata "harmoni/internal/pkg/data"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/mq/publisher"
	commonserver "harmoni/internal/pkg/server"
	"harmoni/internal/pkg/snowflakex"
	commonrepo "harmoni/internal/repository"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
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
		publisher.ProviderSetPublisher,
		commondata.ProviderSet,
		commonrepo.ProviderSet,
		snowflakex.NewSnowflakeNode,
		// etcdx.NewETCDClient,
		repository.ProviderSet,
		usecase.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		commonserver.NewMQServer,
		newApp,
	))
}
