package user

import (
	"context"
	accountentity "harmoni/internal/entity/account"
	"harmoni/internal/entity/paginator"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      int64  `json:"user_id,string" gorm:"not null;uniqueIndex"`
	Name        string `json:"name" gorm:"not null;type:varchar(20)"`
	Email       string `json:"email" gorm:"not null;uniqueIndex;type:varchar(100)"`
	Password    string `json:"-" gorm:"not null;type:varchar(255)"`
	FollowCount int64  `gorm:"not null;default:0"`
	LikeCount   int64  `gorm:"not null;default:0"`
	Avatar      int64  `json:"avatar" gorm:"type:varchar(255)"`
}

type UserBasicInfo struct {
	UserID int64  `json:"user_id,string,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type UserDetail struct {
	UserBasicInfo
	FollowCount int64 `json:"follow_count"`
	LikeCount   int64 `json:"like_count"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) ToBasicInfo(avatarLink string) UserBasicInfo {
	return UserBasicInfo{
		UserID: u.UserID,
		Name:   u.Name,
		Avatar: avatarLink,
	}
}

func ConvertUserToDetailDisplay(u *User, avatarLink string) UserDetail {
	return UserDetail{
		UserBasicInfo: u.ToBasicInfo(avatarLink),
		FollowCount:   u.FollowCount,
		LikeCount:     u.LikeCount,
	}
}

type (
	ModifyStatus uint8
	VerifyType   uint8
)

const (
	NotVerifiedEmailOrPhone ModifyStatus = iota
	VerifiedEmail
	VerifiedPhone
)

const (
	VerifyByEmail VerifyType = iota + 1
	VerifyByPhone
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, bool, error)
	GetByUserID(ctx context.Context, userID int64) (*User, bool, error)
	GetByUserIDs(ctx context.Context, userID []int64) ([]User, error)
	GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[User], error)

	// To ensure that the status corresponds to the correct security information
	// the verification method (email or phone) needs to be specified.
	// you can leave method empty in the implementation
	GetModifyStaus(ctx context.Context, userID int64, verifyType VerifyType, actionType accountentity.AccountActionType) (ModifyStatus, error)
	SetModifyStatus(ctx context.Context, userID int64, status ModifyStatus, verifyType VerifyType, actionType accountentity.AccountActionType, statusKeepTime time.Duration) error

	ModifyPassword(ctx context.Context, user *User) error
	ModifyEmail(ctx context.Context, user *User) error

	GetLikeCount(ctx context.Context, userID int64) (int64, bool, error)
	UpdateLikeCount(ctx context.Context, userID int64, likeCount int64) error

	SetAvatarID(ctx context.Context, userID int64, fileID int64) error
	GetAvatarID(ctx context.Context, userID int64) (int64, error)
}
