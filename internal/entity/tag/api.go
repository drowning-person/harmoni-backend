package tag

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
)

type GetTagsRequest struct {
	entity.PageCond
	Order string `query:"order"`
}

type GetTagsReply struct {
	paginator.Page[TagInfo] `json:"tags"`
}

type GetTagDetailRequest struct {
	TagID int64 `params:"id"`
}

type GetTagDetailReply struct {
	TagDetail
}

type CreateTagRequest struct {
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

type CreateTagReply struct {
	TagDetail
}
