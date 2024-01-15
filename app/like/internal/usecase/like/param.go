package like

import (
	entitylike "harmoni/app/like/internal/entity/like"
)

type LikeRequest struct {
	UserID         int64
	TargetUserID   int64
	LikeType       entitylike.LikeType
	TargetObjectID int64
	IsCancel       bool
}
