// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	remind4 "harmoni/app/notification/internal/handler/remind"
	"harmoni/app/notification/internal/infrastructure/conf"
	"harmoni/app/notification/internal/infrastructure/data"
	"harmoni/app/notification/internal/infrastructure/rpc"
	"harmoni/app/notification/internal/pkg/middleware"
	"harmoni/app/notification/internal/repository/notifyconfig"
	"harmoni/app/notification/internal/repository/remind"
	"harmoni/app/notification/internal/server/http"
	"harmoni/app/notification/internal/server/mq"
	remind3 "harmoni/app/notification/internal/service/remind"
	remind2 "harmoni/app/notification/internal/usecase/remind"
	"harmoni/app/notification/internal/usecase/remind/events"
	conf2 "harmoni/internal/conf"
	data2 "harmoni/internal/pkg/data"
	"harmoni/internal/pkg/etcdx"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/server"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/repository"
)

// Injectors from wire.go:

func initApplication(conf3 *conf.Conf, appConf *conf2.App, dataConf *conf2.DB, etcdConf *conf2.ETCD, serverConf *conf2.Server, logConf *conf2.Log, mqConf *conf2.MessageQueue) (*kratos.App, func(), error) {
	zapLogger, err := logger.NewZapLogger(logConf)
	if err != nil {
		return nil, nil, err
	}
	db, cleanup, err := data.NewDB(dataConf, zapLogger)
	if err != nil {
		return nil, nil, err
	}
	dataDB := data2.NewDB(db)
	loggerLogger := logger.NewLogger(zapLogger)
	notifyConfigRepo := notifyconfig.NewNotifyConfigRepo(dataDB, loggerLogger)
	client, err := etcdx.NewETCDClient(etcdConf, zapLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userClient, err := rpc.NewUsergRPC(client)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	node, err := snowflakex.NewSnowflakeNode(appConf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	uniqueIDRepo := repository.NewUniqueIDRepo(node)
	remindRepo, err := remind.New(userClient, dataDB, uniqueIDRepo, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	remindUsecase := remind2.NewRemindUsecase(notifyConfigRepo, remindRepo, dataDB, loggerLogger)
	remindEventsHandler := events.NewLikeEventsHandler(remindUsecase)
	router, err := mq.NewMQRouter(mqConf, remindEventsHandler, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	mqServer := server.NewMQServer(router)
	authUserMiddleware := middleware.NewAuthUserMiddleware(userClient, loggerLogger)
	remindService := remind3.NewRemindService(remindUsecase, loggerLogger)
	remindHandler := remind4.NewRemindHandler(remindService, loggerLogger)
	notificationAPIRouter := http.NewNotificationAPIRouter(remindHandler, loggerLogger)
	ginServer := http.NewHTTPServer(serverConf, loggerLogger, authUserMiddleware, notificationAPIRouter)
	app := newApplication(conf3, mqServer, ginServer, zapLogger, loggerLogger)
	return app, func() {
		cleanup()
	}, nil
}
