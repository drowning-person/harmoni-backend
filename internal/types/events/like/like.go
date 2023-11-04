package like

import "time"

type LikeType uint8

const (
	LikePost LikeType = iota + 1
	LikeComment
	LikeUser
)

type BaseMessage struct {
	LikeType LikeType `json:"like_type,omitempty"`
}

type LikeCreatedMessage struct {
	BaseMessage
	UserID       int64      `json:"user_id,omitempty"`
	TargetUserID int64      `json:"target_user_id,omitempty"`
	LikingID     int64      `json:"liking_id,omitempty"`
	IsCancel     bool       `json:"is_cancel,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

type LikeStoreMessage struct {
	BaseMessage
	// key is liking id, value is like count
	Counts map[int64]int64
}
