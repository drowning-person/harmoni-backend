package notification

import (
	"time"
)

type RemindParticipant struct {
	ID        int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;NOT NULL"`
	RemindID  int64     `gorm:"column:remind_id;type:bigint(20);NOT NULL"`           // 关联的提醒ID
	RpID      int64     `gorm:"column:rp_id;type:bigint(20);NOT NULL"`               // 提醒参与ID
	SenderID  int64     `gorm:"column:sender_id;type:bigint(20);default:0;NOT NULL"` // 发送人id
}

func (RemindParticipant) TableName() string {
	return "remind_participant"
}
