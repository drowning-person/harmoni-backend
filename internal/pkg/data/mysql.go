package data

import (
	"errors"
	"harmoni/internal/conf"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func ReturnErr(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorx.NotFound(reason.DataNotFoundError)
	case err != nil:
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		return nil
	}
}

type ScopeFunc func(*gorm.DB) *gorm.DB

func zapLevelToGORMLevel(zapLevel zapcore.Level) logger.LogLevel {
	switch zapLevel {
	case zap.InfoLevel:
		return logger.Info
	case zap.ErrorLevel:
		return logger.Error
	case zap.WarnLevel:
		return logger.Warn
	default:
		return logger.Info
	}
}

func NewGormDB(conf *conf.DB, logger *zap.Logger) (*gorm.DB, func(), error) {
	l := zapgorm2.New(logger)
	l.LogLevel = zapLevelToGORMLevel(logger.Level())
	l.SetAsDefault()
	db, err := gorm.Open(mysql.Open(conf.Source), &gorm.Config{Logger: l})
	if err != nil {
		return nil, func() {}, err
	}

	sqlDB, err := db.DB()
	cleanFunc := func() {
		sqlDB.Close()
	}
	if err != nil {
		return nil, cleanFunc, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(int(conf.MaxIdleConn))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(int(conf.MaxOpenConn))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime.AsDuration())

	return db, cleanFunc, nil
}
