package notifyconfig

import (
	"context"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"
)

type NotifyConfig struct {
	ID              int64
	MessageTemplate string            // 消息模版
	Action          action.Action     // 动作类型
	ObjectType      object.ObjectType // 被操作对象类型
	NotifyChannel   string            // 为某个通知类型设置一个或多个推送渠道
}

type NotifyConfigRepository interface {
	Get(ctx context.Context, action action.Action, objectType object.ObjectType) (*NotifyConfig, error)
}
