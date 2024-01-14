package v1

import (
	"time"

	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/internal/pkg/httpx"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"
)

type UnReadRequest struct {
	UserID int64         `json:"user_id"`
	Action action.Action `form:"action"`
}

type UnReadResponse struct {
	Count int64 `json:"count"`
}

type ListRemindRequest struct {
	httpx.Page
	UserID int64         `json:"-"`
	Action action.Action `json:"-"`
}

type Remind struct {
	RemindID   int64
	Senders    v1.UserBasicList  // 发送方 只显示前4个
	Recipient  *v1.UserBasic     // 接受通知的用户的ID
	ObjectID   int64             // 目标对象ID
	ObjectType object.ObjectType // 被操作对象类型
	Content    string            // 通知内容
}

type ListRemindResponse struct {
	httpx.PageResp
	Reminds []*Remind `json:"reminds"`
}

type LikeRemindDetailRequest struct {
	httpx.Page
	RemindID int64         `json:"-"`
	UserID   int64         `json:"-"`
	Action   action.Action `json:"-"`
}

type LikeRemindDetailObject struct {
	ObjectID   int64             `json:"oid"`   // 目标对象ID
	ObjectType object.ObjectType `json:"type"`  // 被操作对象类型
	Title      string            `json:"title"` // 对象标题
	Describe   string            `json:"desc"`  // 对象描述信息
	CreatedAt  *time.Time        `json:"ctime"` // 创建时间
}

type LikeRemindDetailItem struct {
	User      *v1.UserBasic `json:"user"`
	CreatedAt *time.Time    `json:"ctime"`
}

type LikeRemindDetailResponse struct {
	httpx.PageResp
	// TODO wait for object service
	Object *LikeRemindDetailObject `json:"object"`
	Items  []*LikeRemindDetailItem `json:"items"`
}
