package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   int64  `json:"user_id,string"`
	Name     string `json:"name" gorm:"not null;type:varchar(20)"`
	Email    string `json:"email" gorm:"uniqueIndex;type:varchar(100)"`
	Password string `json:"-" gorm:"not null;type:varchar(255)"`
}

type UserDetail struct {
	UserID    int64     `json:"user_id,string"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_time"`
}

type UserInfo struct {
	UserID int64  `json:"user_id,string"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func IsUserExist(email string) (bool, error) {
	var count int64
	if err := DB.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}
