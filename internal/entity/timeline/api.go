package timeline

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
)

type GetUserTimeLineRequest struct {
	entity.PageCond
	AuthorID int64 `json:"author_id" validate:"required"`
	UserID   int64 `json:"-"`
}

type GetUserTimeLineReply struct {
	paginator.Page[postentity.PostBasicInfo]
}

type GetHomeTimeLineRequest struct {
	entity.PageCond
	UserID int64 `json:"-"`
}

type GetHomeTimeLineReply struct {
	paginator.Page[postentity.PostBasicInfo]
}
