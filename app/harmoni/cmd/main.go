package main

import (
	"encoding/json"
	"harmoni/app/harmoni/internal/cron"
	"harmoni/app/harmoni/internal/infrastructure/config"

	"harmoni/app/harmoni/internal/server/http"
	"harmoni/app/harmoni/internal/server/mq"
	"harmoni/internal/pkg/validator"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	// go build -ldflags "-X main.Version=x.y.z"
	Version string
)

func newApplication(
	conf *config.Config,
	etcdclient *etcdclient.Client,
	fiberServer *http.FiberServer,
	grpcServer *grpc.Server,
	cronApp *cron.ScheduledTaskManager,
	mqServer *mq.MQServer,
	logger *zap.SugaredLogger,
	kratosLogger log.Logger,
) *kratos.App {
	data, err := json.Marshal(conf)
	if err != nil {
		logger.Warnf("print config failed: %s", data)
	} else {
		logger.Infof("config json: %+v", string(data))
	}
	return kratos.New(
		kratos.Name(conf.App.ServiceName),
		kratos.Version(Version),
		kratos.Server(
			fiberServer,
			grpcServer,
			cronApp,
			mqServer,
		),
		kratos.Logger(kratosLogger),
		kratos.Registrar(etcd.New(etcdclient)),
	)
}

func main() {
	cfg, err := config.ReadConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}
	app, clean, err := initApplication(
		cfg,
		cfg.App,
		cfg.DB,
		cfg.Redis,
		cfg.Auth,
		cfg.Email,
		cfg.Like,
		cfg.MessageQueue,
		cfg.FileStorage,
		cfg.ETCD,
		cfg.Server,
		cfg.Log)
	defer clean()
	if err != nil {
		panic(err)
	}
	err = validator.InitTrans(cfg.App.Locale)
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
