package entity

import (
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	DefaultRedisValue = "*" //redis中key对应的预设值，防脏读
)

type Int64Slice []int64

// MarshalJSON 实现 json.Marshaler 接口
func (is Int64Slice) MarshalJSON() ([]byte, error) {
	strSlice := make([]string, len(is))

	for i, num := range is {
		strSlice[i] = strconv.FormatInt(num, 10)
	}

	return json.Marshal(strSlice)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (is *Int64Slice) UnmarshalJSON(data []byte) error {
	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err != nil {
		return err
	}

	intSlice := make(Int64Slice, len(strSlice))

	for i, str := range strSlice {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}

		intSlice[i] = num
	}

	*is = intSlice

	return nil
}

type BaseModel struct {
	ID uint64 `gorm:"primarykey"`
	TimeMixin
	SoftDeletes
}

type TimeMixin struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDeletes struct {
	DeletedAt gorm.DeletedAt
}
