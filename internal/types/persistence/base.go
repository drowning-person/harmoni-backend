package persistence

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID uint64 `gorm:"primarykey"`
	TimeMixin
}

type BaseModelWithSoftDelete struct {
	BaseModel
	SoftDelete
}

type BaseModelWithSoftDeleteUnix struct {
	BaseModel
	soft_delete.DeletedAt
}

type TimeMixin struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt
}
