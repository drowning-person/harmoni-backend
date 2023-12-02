package notification

import (
	"harmoni/app/notification/internal/entity/notifyconfig"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"
	"time"
)

type NotifyConfig struct {
	ID              int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;NOT NULL"`
	MessageTemplate string    `gorm:"column:message_template;type:varchar(1024);NOT NULL"` // 消息模版
	Action          int8      `gorm:"column:action;type:tinyint(4);NOT NULL"`              // 动作类型
	ObjectType      int8      `gorm:"column:object_type;type:tinyint(4);NOT NULL"`         // 被操作对象类型
	NotifyChannel   string    `gorm:"column:notify_channel;type:varchar(255);NOT NULL"`    // 为某个通知类型设置一个或多个推送渠道
}

func (NotifyConfig) TableName() string {
	return "notify_config"
}

func (n *NotifyConfig) ToDomain() *notifyconfig.NotifyConfig {
	return &notifyconfig.NotifyConfig{
		ID:              n.ID,
		MessageTemplate: n.MessageTemplate,
		Action:          action.Action(n.Action),
		ObjectType:      object.ObjectType(n.ObjectType),
		NotifyChannel:   n.NotifyChannel,
	}
}
