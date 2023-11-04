package config

import (
	"fmt"
	"harmoni/internal/pkg/filesystem/policy"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App          *App          `mapstructure:"app"`
	Auth         *Auth         `mapstructure:"auth"`
	DB           *DB           `mapstructure:"db"`
	Log          *Log          `mapstructure:"log"`
	Redis        *Redis        `mapstructure:"redis"`
	Email        *Email        `mapstructure:"email"`
	MessageQueue *MessageQueue `mapstructure:"messageQueue"`
	FileStorage  *FileStorage  `mapstructure:"fileStorage"`
}

type App struct {
	Debug     bool   `default:"false" mapstructure:"debug"`
	Addr      string `default:"127.0.0.1:80" mapstructure:"addr"`
	BaseURL   string `mapstructure:"base_url"`
	StartTime string `mapstructure:"start_time"`
	AppID     int64  `mapstructure:"app_id"`
	Locale    string `mapstructure:"locale"`
}

func SetAppDefault(v *viper.Viper) {
	v.SetDefault("app", map[string]interface{}{
		"debug":      false,
		"addr":       "127.0.0.1:80",
		"base_url":   "localhost",
		"start_time": time.Now().Format("2006-01-02"),
		"app_id":     1,
		"locale":     "zh",
	})
}

func SetAppEnv(v *viper.Viper) error {
	err := v.BindEnv("app.debug", "HARMONI_DEBUG")
	if err != nil {
		return err
	}

	err = v.BindEnv("app.addr", "HARMONI_ADDR")
	if err != nil {
		return err
	}

	return nil
}

type Auth struct {
	TokenExpire        time.Duration `mapstructure:"token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
	Secret             string        `mapstructure:"secret"`
}

func SetAuthDefault(v *viper.Viper) {
	v.SetDefault("auth", map[string]interface{}{
		"token_expire":         "5m",
		"refresh_token_expire": "336h",
		"secret":               "IKNOWWHATIAMDOING",
	})
}

func SetAuthEnv(v *viper.Viper) error {
	err := v.BindEnv("auth.secret", "HARMONI_AUTH_SECRET")
	if err != nil {
		return err
	}

	return nil
}

type DB struct {
	Driver          string
	Source          string
	ConnMaxLifeTime time.Duration `mapstructure:"conn_max_life_time" yaml:"conn_max_life_time,omitempty"`
	MaxOpenConn     int           `mapstructure:"max_open_conn" yaml:"max_open_conn,omitempty" default:"8"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"  yaml:"max_idle_conn,omitempty" default:"8"`
}

func SetDBDefault(v *viper.Viper) {
	v.SetDefault("db", map[string]interface{}{
		"driver":             "mysql",
		"source":             "root:123456@tcp(127.0.0.1:3306)/harmoni?parseTime=True",
		"conn_max_life_time": "1h",
		"max_open_conn":      runtime.NumCPU() * 2,
		"max_idle_conn":      runtime.NumCPU() * 2,
	})
}

func SetDBEnv(v *viper.Viper) error {
	err := v.BindEnv("db.driver", "HARMONI_DB_DRIVER")
	if err != nil {
		return err
	}

	err = v.BindEnv("db.source", "HARMONI_DB_SOURCE")
	if err != nil {
		return err
	}

	return nil
}

type Log struct {
	Level string `default:"info"`
	Path  string `default:"./log/harmoni.log"`
	Color bool
}

func SetLogDefault(v *viper.Viper) {
	v.SetDefault("log", map[string]interface{}{
		"level": "info",
		"path":  "./log/harmoni.log",
		"color": true,
	})
}

type Redis struct {
	IP           string        `mapstructure:"ip"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int8          `mapstructure:"database"`
	PoolSize     int           `mapstructure:"pool_size"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func SetRedisEnv(v *viper.Viper) error {
	err := v.BindEnv("redis.ip", "HARMONI_REDIS_IP")
	if err != nil {
		return err
	}

	err = v.BindEnv("redis.port", "HARMONI_REDIS_PORT")
	if err != nil {
		return err
	}

	err = v.BindEnv("redis.password", "HARMONI_REDIS_PASSWORD")
	if err != nil {
		return err
	}

	return nil
}

func SetRedisDefault(v *viper.Viper) {
	v.SetDefault("redis", map[string]interface{}{
		"ip":            "127.0.0.1",
		"port":          6379,
		"password":      "",
		"database":      0,
		"pool_size":     runtime.NumCPU() * 2,
		"read_timeout":  "3s",
		"write_timeout": "3s",
	})
}

type Email struct {
	Host     string        `mapstructure:"host"`
	Port     string        `mapstructure:"port"`
	UserName string        `mapstructure:"user_name"`
	Password string        `mapstructure:"password"`
	FromName string        `mapstructure:"from_name"`
	CodeTTL  time.Duration `mapstructure:"code_ttl"`
}

func SetEmailDefault(v *viper.Viper) {
	v.SetDefault("email", map[string]interface{}{
		"code_ttl": "5m",
	})
}

func SetEmailEnv(v *viper.Viper) error {
	err := v.BindEnv("email.host", "HARMONI_EMAIL_HOST")
	if err != nil {
		return err
	}

	err = v.BindEnv("email.port", "HARMONI_EMAIL_PORT")
	if err != nil {
		return err
	}

	err = v.BindEnv("email.user_name", "HARMONI_EMAIL_USERNAME")
	if err != nil {
		return err
	}

	err = v.BindEnv("email.password", "HARMONI_EMAIL_PASSWORD")
	if err != nil {
		return err
	}

	return nil
}

type MessageQueue struct {
	RabbitMQ *RabbitMQConf `mapstructure:"rabbitmq,omitempty"`
}

type RabbitMQConf struct {
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Host     string `mapstructure:"host,omitempty"`
	Port     int    `mapstructure:"port,omitempty"`
	VHost    string `mapstructure:"vhost,omitempty"`
}

func (c *RabbitMQConf) BuildURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.VHost)
}

type FilePolicyOption struct {
	// 允许的文件扩展名
	FileType []string `json:"fileType"`
	// 分片上传的分片大小
	ChunkSize uint64 `json:"chunkSize,omitempty"`
}

type FilePolicy struct {
	Type       string `mapstructure:"type"`
	BucketName string `mapstructure:"bucketName"`
	MaxSize    uint64 `mapstructure:"maxSize"`
	// key type, value dir
	DirRule map[policy.FileType]string `mapstructure:"dirRule"`
	Option  FilePolicyOption           `mapstructure:"option"`
}

type FileStorage struct {
	DefaultAvatar        string        `mapstructure:"defaultAvatar"`
	UploadSessionTimeout time.Duration `mapstructure:"uploadSessionTimeout"`
	Policy               *FilePolicy   `mapstructure:"policy,omitempty"`
}

func SetFileStorageDefault(v *viper.Viper) {
	v.SetDefault("file_storage", map[string]interface{}{
		"uploadSessionTimeout": "2h",
		"policy": FilePolicy{
			Type:       "local",
			BucketName: "./static",
			MaxSize:    1024 * 1024 * 100,
			Option: FilePolicyOption{
				ChunkSize: 5 * 1024 * 1024,
			},
		},
	})
}

func ReadConfig(filePath string) (*Config, error) {
	v := viper.New()

	SetAppDefault(v)
	SetAuthDefault(v)
	SetDBDefault(v)
	SetLogDefault(v)
	SetRedisDefault(v)
	SetEmailDefault(v)
	SetFileStorageDefault(v)

	filename := path.Base(filePath)
	fileext := path.Ext(filePath)
	filepre := strings.TrimSuffix(filename, fileext)
	v.SetConfigName(filepre)
	v.SetConfigType(strings.TrimPrefix(fileext, "."))
	v.AddConfigPath(path.Dir(filePath))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = SetAppEnv(v)
	if err != nil {
		return nil, err
	}

	err = SetAuthEnv(v)
	if err != nil {
		return nil, err
	}

	err = SetDBEnv(v)
	if err != nil {
		return nil, err
	}

	err = SetRedisEnv(v)
	if err != nil {
		return nil, err
	}

	err = SetEmailEnv(v)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
