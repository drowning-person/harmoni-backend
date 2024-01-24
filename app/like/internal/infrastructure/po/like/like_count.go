package like

import (
	objectv1 "harmoni/api/common/object/v1"
	"harmoni/internal/types/persistence"
)

type LikeCount struct {
	persistence.BaseModelWithSoftDeleteUnix
	Counts     int64               `gorm:"not null;"`
	ObjectID   int64               `gorm:"not null;"`
	ObjectType objectv1.ObjectType `gorm:"not null;type:TINYINT UNSIGNED"`
}
