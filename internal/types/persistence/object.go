package persistence

import (
	"harmoni/internal/types/object"

	"gorm.io/gorm"
)

func ByObject(objectID int64, objectType object.ObjectType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("object_id = ?", objectID).
			Where("object_type = ?", objectType)
	}
}
