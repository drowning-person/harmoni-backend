package like

import (
	"context"
	"encoding/json"
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/pkg/common"
	"time"
)

type LikeType uint8

const (
	LikePost LikeType = iota + 1
	LikeComment
	LikeUser
)

type Like struct {
	ID uint `gorm:"primarykey;type:BIGINT UNSIGNED not NULL AUTO_INCREMENT;"`
	entity.TimeMixin
	UserID   int64    `gorm:"not null"`
	LikingID int64    `gorm:"not null"`
	LikeType LikeType `gorm:"not null;type:TINYINT UNSIGNED"`
	Canceled bool     `gorm:"not null;default:0;"`
}

func (*Like) TableName() string {
	return "like"
}

type LikeCacheInfo struct {
	LikingID  int64
	UpdatedAt int64 `gorm:"serializer:unixtime"`
}

func (r *LikeCacheInfo) ToJSONString() string {
	codeBytes, _ := json.Marshal(r)
	return common.BytesToString(codeBytes)
}

func (r *LikeCacheInfo) FromJSONString(data string) error {
	return json.Unmarshal(common.StringToBytes(data), r)
}

type LikeMessageType uint8

const (
	LikeCountMessage LikeMessageType = iota + 1
	LikeActionMessage
)

type CountMessage struct {
	// key is liking id, value is like count
	Counts map[int64]int64
}

type ActionMessage struct {
	UserID    int64      `json:"user_id,omitempty"`
	LikingID  int64      `json:"liking_id,omitempty"`
	IsCancel  bool       `json:"is_cancel,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

var (
	LikeTypeList = []LikeType{LikeUser, LikePost, LikeComment}
)

type LikeMessage struct {
	Type          LikeMessageType `json:"type,omitempty"`
	LikeType      LikeType        `json:"like_type,omitempty"`
	CountMessage  *CountMessage   `json:"count_message,omitempty"`
	ActionMessage *ActionMessage  `json:"action_message,omitempty"`
}

type LikeRepository interface {
	Like(ctx context.Context, like *Like, targetUserID int64, isCancel bool) error
	Save(ctx context.Context, like *Like, isCancel bool) error
	LikeCount(ctx context.Context, like *Like) (int64, bool, error)
	BatchLikeCount(ctx context.Context, likeType LikeType) (map[int64]int64, error)
	BatchLikeCountByIDs(ctx context.Context, likingIDs []int64, likeType LikeType) (map[int64]int64, error)
	// UpdateLikeCount(ctx context.Context, like *Like, count int8) error
	ListLikingIDs(ctx context.Context, query *LikeQuery) (paginator.Page[int64], error)
	IsLiking(ctx context.Context, like *Like) (bool, error)
	CacheLikeCount(ctx context.Context, like *Like, count int64) error
}
