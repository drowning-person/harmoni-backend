package notification

import (
	"time"
)

type NotifyRemind struct {
	ID           int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;NOT NULL"`
	RemindID     int64     `gorm:"column:remind_id;type:bigint(20);NOT NULL"`
	RecipientID  int64     `gorm:"column:recipient_id;type:bigint(20);NOT NULL"`          // 接受通知的用户的ID
	Action       int8      `gorm:"column:action;type:tinyint(4);default:0;NOT NULL"`      // 动作类型
	ObjectID     int64     `gorm:"column:object_id;type:bigint(20);default:0;NOT NULL"`   // 目标对象ID
	ObjectType   int8      `gorm:"column:object_type;type:tinyint(4);default:0;NOT NULL"` // 被操作对象类型
	Content      string    `gorm:"column:content;type:varchar(1024);default:0;NOT NULL"`  // 通知内容
	LastReadTime time.Time `gorm:"column:last_read_time;type:datetime;NOT NULL"`          // 最近已读时间
}

func (NotifyRemind) TableName() string {
	return "notify_remind"
}
