package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"harmoni/app/notification/internal/infrastructure/conf"
	"harmoni/internal/pkg/server"
	"harmoni/internal/pkg/validator"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/transport/gin"
	"go.uber.org/zap"
)

var (
	// go build -ldflags "-X main.Version=x.y.z"
	Version string

	ConfigPath = "./config/config.yaml"
)

func init() {
	flag.StringVar(&ConfigPath, "conf", "./config/config.yaml", "config path, eg: -conf config.yaml")
}

func newApplication(
	conf *conf.Conf,
	mqServer *server.MQServer,
	httpServer *gin.Server,
	logger *zap.Logger,
	kratosLogger log.Logger,
) *kratos.App {
	data, err := json.Marshal(conf)
	if err != nil {
		logger.Sugar().Warnf("print config failed: %s", data)
	} else {
		logger.Sugar().Infof("config json: %+v", string(data))
	}
	return kratos.New(
		kratos.Name(conf.App.ServerName),
		kratos.Version(Version),
		kratos.Server(
			httpServer, mqServer,
		),
		kratos.Logger(kratosLogger),
	)
}

func main() {
	flag.Parse()
	cfg := config.New(
		config.WithSource(
			file.NewSource(ConfigPath),
		))
	err := cfg.Load()
	if err != nil {
		panic(err)
	}
	conf := conf.Conf{}
	err = cfg.Scan(&conf)
	if err != nil {
		panic(err)
	}
	confJson, err := json.Marshal(&conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", confJson)
	app, clean, err := initApplication(
		&conf,
		conf.GetApp(),
		conf.GetDb(),
		conf.GetEtcd(),
		conf.GetServer(),
		conf.GetLog(),
		conf.GetMessageQueue())
	defer clean()
	if err != nil {
		panic(err)
	}
	err = validator.InitTrans(conf.App.GetLocale())
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
