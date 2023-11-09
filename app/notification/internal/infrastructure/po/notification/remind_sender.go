package notification

import "harmoni/internal/types/persistence"

type RemindSender struct {
	persistence.BaseModel
	RemindID int64 `gorm:"column:remind_id;type:bigint;NOT NULL"`           // 关联的提醒ID
	RPID     int64 `gorm:"column:rp_id;type:bigint;NOT NULL"`               // 提醒参与ID
	SenderID int64 `gorm:"column:sender_id;type:bigint;default:0;NOT NULL"` // 发送人id
}

func (m *RemindSender) TableName() string {
	return "remind_sender"
}
