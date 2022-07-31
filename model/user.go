package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   int64  `json:"user_id,string"`
	Name     string `json:"name" gorm:"not null;type:varchar(255)"`
	Email    string `json:"email" gorm:"uniqueIndex;type:varchar(100)"`
	Password string `json:"-" gorm:"not null;type:varchar(255)"`
}

func (u *User) HashAndSalt(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsExist() (bool, error) {
	if err := DB.Where("email = ?", u.Email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
