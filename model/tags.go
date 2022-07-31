package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	TagID        int64
	TagName      string `gorm:"varchar(128);unique"`
	Introduction string `gorm:"varchar(256);unique"`
}
