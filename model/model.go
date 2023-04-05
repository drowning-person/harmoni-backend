package model

import (
	"harmoni/config"
	"harmoni/pkg/zap"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var DB *gorm.DB

func InitMysql(conf *config.DB) error {
	logger := zapgorm2.New(zap.Logger)
	logger.SetAsDefault()
	db, err := gorm.Open(mysql.Open(conf.Source), &gorm.Config{Logger: logger})
	DB = db
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(12)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(12)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute)

	DB.AutoMigrate(&User{}, &Tag{}, &Post{}, &Comment{})

	return nil
}
