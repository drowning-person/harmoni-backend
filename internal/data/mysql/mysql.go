package mysql

import (
	"harmoni/internal/conf"
	commententity "harmoni/internal/entity/comment"
	postentity "harmoni/internal/entity/post"
	tagentity "harmoni/internal/entity/tag"
	userentity "harmoni/internal/entity/user"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func NewDB(conf *conf.DB, logger *zap.Logger) (*gorm.DB, error) {
	l := zapgorm2.New(logger)
	l.SetAsDefault()
	db, err := gorm.Open(mysql.Open(conf.Source), &gorm.Config{Logger: l})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	logger.Sugar().Debugf("db conf is %#v", *conf)

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	db.AutoMigrate(&userentity.User{}, &commententity.Comment{}, &postentity.Post{}, &tagentity.Tag{})

	return db, nil
}
