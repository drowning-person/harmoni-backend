package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"harmoni/app/like/internal/conf"
	"harmoni/internal/pkg/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(
	logger log.Logger,
	gs *grpc.Server,
	ms *server.MQServer,
) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			ms,
		),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	conf := conf.Conf{}
	if err := c.Scan(&conf); err != nil {
		panic(err)
	}
	confJson, err := json.Marshal(&conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", confJson)
	app, cleanup, err := wireApp(
		&conf,
		conf.GetApp(),
		conf.GetDb(),
		conf.GetEtcd(),
		conf.GetServer(),
		conf.GetLog(),
		conf.GetMessageQueue(),
		conf.GetRedis(),
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
