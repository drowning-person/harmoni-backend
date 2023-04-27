package comment

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
)

type CreateCommentRequest struct {
	UserID int64 `json:"-"`
	// object id
	ObjectID string `json:"oid" validate:"required" label:"对象ID"`
	ParentID int64  `json:"pid"`
	RootID   int64  `json:"rid"`
	Content  string `json:"content" validate:"required,gte=10,lte=512" label:"评论内容"`
}

type CreateCommentReply struct {
	CommentDetail
}

type GetCommentsRequest struct {
	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
	// object id
	ObjectID string `query:"oid" validate:"required" label:"对象ID"`
	// root id
	RootID string `query:"rid" label:"根ID"`
	// query condition
	QueryCond string `query:"cond" validate:"omitempty,oneof=newest" label:"排序"`
}

type GetCommentsReply struct {
	paginator.Page[CommentDetail]
}

func ConvertPageReqToCommentQuery(req *GetCommentsRequest) CommentQuery {
	return CommentQuery{
		PageCond: entity.PageCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		ObjectID:  req.ObjectID,
		RootID:    req.RootID,
		QueryCond: req.QueryCond,
	}
}
