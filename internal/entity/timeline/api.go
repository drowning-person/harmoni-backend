package timeline

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
)

type GetUserTimeLineRequest struct {
	entity.PageCond
	UserID int64 `json:"user_id,omitempty"`
}

type GetUserTimeLineReply struct {
	paginator.Page[postentity.PostDetail]
}

type GetHomeTimeLineRequest struct {
	entity.PageCond
	UserID int64 `json:"-"`
}

type GetHomeTimeLineReply struct {
	paginator.Page[postentity.PostDetail]
}
