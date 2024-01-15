package like

import (
	"harmoni/app/like/internal/entity/like"
	"harmoni/internal/types/persistence"
)

type Like struct {
	persistence.BaseModelWithSoftDeleteUnix
	UserID       int64         `gorm:"not null"`
	TargetUserID int64         `gorm:"not null"`
	LikingID     int64         `gorm:"not null;uniqueIndex"`
	LikeType     like.LikeType `gorm:"not null;type:TINYINT UNSIGNED"`
	ObjectID     int64         `gorm:"not null;"`
}

func (Like) TableName() string {
	return "like"
}

func FromDomain(like *like.Like) *Like {
	return &Like{
		UserID:       like.User.GetId(),
		TargetUserID: like.TargetUser.GetId(),
		LikingID:     like.LikingID,
		LikeType:     like.LikeType,
		ObjectID:     like.ObjectID,
	}
}
