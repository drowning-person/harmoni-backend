package follow

import (
	"context"
	"harmoni/app/harmoni/internal/entity/paginator"
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	FollowerID   int64        `gorm:"not null;index"`                 // 关注者用户 ID
	FollowingID  int64        `gorm:"not null;index"`                 // 被关注对象 ID
	FollowedType FollowedType `gorm:"not null;type:TINYINT UNSIGNED"` // 被关注对象类型（'user' 或 'topic'）
}

func (*Follow) TableName() string {
	return "follow"
}

type FollowBasicInfo struct {
	FollowerID   int64        `json:"follower_id,omitempty,string"`  // 关注者用户 ID
	FollowingID  int64        `json:"following_id,omitempty,string"` // 被关注对象 ID
	FollowedType FollowedType `json:"type,omitempty"`                // 被关注对象类型（'user' 或 'topic'）
	CreatedAt    time.Time    `json:"created_at,omitempty"`
}

func ConvertFollowToBasic(f *Follow) *FollowBasicInfo {
	return &FollowBasicInfo{
		FollowerID:   f.FollowerID,
		FollowingID:  f.FollowingID,
		FollowedType: f.FollowedType,
		CreatedAt:    f.CreatedAt,
	}
}

type FollowedType uint8

const (
	FollowUser FollowedType = iota + 1
	FollowTag
)

type FollowRepository interface {
	Follow(ctx context.Context, follow *Follow) error
	FollowCancel(ctx context.Context, follow *Follow) error
	GetFollowersPage(ctx context.Context, followQuery *FollowQuery) (paginator.Page[int64], error)
	GetFollowingsPage(ctx context.Context, followQuery *FollowQuery) (paginator.Page[int64], error)
	GetFollowingUsersAll(ctx context.Context, userID int64) ([]int64, error)
	IsFollowing(ctx context.Context, follow *Follow) (bool, error)
	AreFollowEachOther(ctx context.Context, userIDx int64, userIDy int64) (bool, error)
}
