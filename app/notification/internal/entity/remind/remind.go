package remind

import (
	"context"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/notification/internal/entity/notifyconfig"
	"harmoni/internal/pkg/paginator"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"
	"time"
)

type Remind struct {
	RemindID     int64
	Recipient    *v1.UserBasic     // 接受通知的用户的ID
	Senders      []*v1.UserBasic   // 发送人
	Action       action.Action     // 动作类型
	ObjectID     int64             // 目标对象ID
	ObjectType   object.ObjectType // 被操作对象类型
	Content      string            // 通知内容
	LastReadTime time.Time         // 最近已读时间
}

func (r *Remind) BuildContent(config *notifyconfig.NotifyConfig) {
	if config == nil {
		r.Content = "提醒"
		return
	}
	r.Content = config.MessageTemplate
}

type ListReq struct {
	*paginator.PageRequest
	UserID      int64
	Action      action.Action
	SenderCount int
}

type CreateReq struct {
	RecipientID  int64
	SenderIDs    []int64
	Action       action.Action
	ObjectID     int64
	ObjectType   object.ObjectType
	Content      string
	LastReadTime *time.Time
}

type UpdateLastReadTimeReq struct {
	UserID   int64
	Action   action.Action
	ReadTime time.Time
}

type CountReq struct {
	UserID int64
	Action action.Action
	UnRead bool
}

type ListRemindSendersReq struct {
	*paginator.PageRequest
	RemindID int64
	Action   action.Action
}

type RemindSender struct {
	Sender    *v1.UserBasic // 发送人
	CreatedAt *time.Time    // 创建时间
}

type RemindRepository interface {
	Create(ctx context.Context, req *CreateReq) error
	List(ctx context.Context, req *ListReq) (*paginator.Page[*Remind], error)
	Count(ctx context.Context, req *CountReq) (int64, error)
	UpdateLastReadTime(ctx context.Context, req *UpdateLastReadTimeReq) error
	ListRemindSenders(ctx context.Context, req *ListRemindSendersReq) (*paginator.Page[*RemindSender], error)
}
