package config

import (
	"time"

	cfg "github.com/gookit/config/v2"
	gyaml "github.com/gookit/config/v2/yaml"
)

type Config struct {
	App   *App
	Auth  *Auth
	DB    *DB
	Log   *Log
	Redis *Redis
}

type App struct {
	Addr string `default:"127.0.0.1:80"`
}

type Auth struct {
	TokenExpire int64 `default:"1800"`
	Secret      string
}

type DB struct {
	Driver string
	Source string
}

type Log struct {
	Level string `default:"info"`
	Path  string `default:"./log/harmoni.log"`
}

type Redis struct {
	IP           string
	Port         int
	Password     string
	Database     int8
	PoolSize     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func ReadConfig(filePath string) (*Config, error) {
	c := &Config{}
	cfg.WithOptions(cfg.ParseEnv, cfg.ParseDefault, cfg.ParseTime)
	cfg.WithOptions(func(opt *cfg.Options) {
		opt.DecoderConfig.TagName = "config"
	})

	// add driver for support yaml content
	cfg.AddDriver(gyaml.Driver)

	err := cfg.LoadFiles(filePath)
	if err != nil {
		panic(err)
	}

	err = cfg.Decode(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
