package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	DefaultRedisValue = "*" //redis中key对应的预设值，防脏读
)

type TimeMixin struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDeletes struct {
	DeletedAt gorm.DeletedAt
}
