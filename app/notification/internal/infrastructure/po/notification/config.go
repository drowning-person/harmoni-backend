package notification

import "harmoni/internal/types/persistence"

type NotifyConfig struct {
	persistence.BaseModel
	MessageTemplate string `gorm:"column:message_template;type:varchar(1024);NOT NULL"` // 消息模版
	Action          int8   `gorm:"column:action;type:tinyint;NOT NULL"`                 // 动作类型
	ObjectType      int8   `gorm:"column:object_type;type:tinyint;NOT NULL"`            // 被操作对象类型
	NotifyChannel   string `gorm:"column:notify_channel;type:varchar(255);NOT NULL"`    // 为某个通知类型设置一个或多个推送渠道
}

func (m *NotifyConfig) TableName() string {
	return "notify_config"
}
