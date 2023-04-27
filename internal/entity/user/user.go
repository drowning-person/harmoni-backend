package user

import (
	"context"
	"harmoni/internal/entity/paginator"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   int64  `json:"user_id,string" gorm:"not null;uniqueIndex"`
	Name     string `json:"name" gorm:"not null;type:varchar(20)"`
	Email    string `json:"email" gorm:"uniqueIndex;type:varchar(100)"`
	Password string `json:"-" gorm:"not null;type:varchar(255)"`
}

type BasicUserInfo struct {
	UserID int64  `json:"user_id,string"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UserDetail struct {
	BasicUserInfo
}

func ConvertUserToDisplay(u *User) BasicUserInfo {
	return BasicUserInfo{
		UserID: u.UserID,
		Name:   u.Name,
		Email:  u.Email,
	}
}

func ConvertUserToDetailDisplay(u *User) UserDetail {
	return UserDetail{
		BasicUserInfo: BasicUserInfo{
			UserID: u.UserID,
			Name:   u.Name,
			Email:  u.Email,
		},
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, bool, error)
	GetByUserID(ctx context.Context, userID int64) (*User, bool, error)
	GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[User], error)
}
