package like

import (
	v1 "harmoni/api/common/object/v1"
)

type LikeRequest struct {
	UserID         int64
	TargetUserID   int64
	ObjectType     v1.ObjectType
	TargetObjectID int64
	IsCancel       bool
}
