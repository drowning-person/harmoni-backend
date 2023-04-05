package zap

import (
	"harmoni/config"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InitLogger(conf *config.Log) error {
	dirname, filename := path.Split(conf.Path)
	exist, err := PathExists(dirname)
	if err != nil {
		return err
	}
	if !exist {
		err := os.Mkdir(dirname, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(dirname+filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	level := zap.InfoLevel
	logLevel := strings.ToLower(conf.Level)
	switch logLevel {
	case "info":
		level = zap.InfoLevel
	case "debug":
		level = zap.DebugLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	}

	Logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(file), level))
	Logger = Logger.WithOptions(zap.AddCaller(), zap.Development())

	zap.ReplaceGlobals(Logger)
	return nil
}
