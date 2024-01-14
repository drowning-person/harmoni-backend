package like

import (
	"harmoni/app/like/internal/entity/like"
	"harmoni/internal/types/persistence"
)

type Like struct {
	persistence.BaseModelWithSoftDelete
	UserID       int64         `gorm:"not null"`
	TargetUserID int64         `gorm:"not null"`
	LikingID     int64         `gorm:"not null;uniqueIndex"`
	LikeType     like.LikeType `gorm:"not null;type:TINYINT UNSIGNED"`
	Canceled     bool          `gorm:"not null;default:0;"`
}

func (Like) TableName() string {
	return "like"
}
