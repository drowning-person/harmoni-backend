package logger

import (
	"harmoni/internal/conf"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func NewZapLogger(conf *conf.Log) (*zap.Logger, error) {
	dirname, filename := path.Split(conf.Path)
	exist, err := PathExists(dirname)
	if err != nil {
		return nil, err
	}
	if !exist {
		err := os.Mkdir(dirname, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(dirname+filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
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

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(file), level)
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level)
	tee := zapcore.NewTee(consoleCore, core)
	logger := zap.New(tee, zap.AddCaller(), zap.Development())

	zap.ReplaceGlobals(logger)

	return logger, nil
}
