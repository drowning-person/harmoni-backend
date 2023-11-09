package persistence

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID uint64 `gorm:"primarykey"`
	TimeMixin
}

type BaseModelWithSoftDelete struct {
	BaseModel
	SoftDelete
}

type TimeMixin struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt
}
