package post

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
)

type GetPostsRequest struct {
	entity.PageCond
	Order string `query:"order" validate:"omitempty,oneof=newest score" label:"排序"`
}

type GetPostsReply struct {
	paginator.Page[PostDetail] `json:"posts"`
}

type GetPostDetailRequest struct {
	PostID int64 `params:"id" label:"帖子ID"`
}

type GetPostDetailReply struct {
	PostDetail
}

type CreatePostRequest struct {
	TagIDs  entity.Int64Slice `json:"tag_ids" validate:"lte=4" label:"话题ID"`
	UserID  int64             `json:"-"`
	Title   string            `json:"title" validate:"required,gte=3,lte=128" label:"帖子标题"`
	Content string            `json:"content" validate:"required,gte=10,lte=512" label:"帖子内容"`
}

type CreatePostReply struct {
	PostDetail
}

type LikePostRequest struct {
	PostID int64 `json:"post_id,string" validate:"required" label:"帖子ID"`
	UserID int64 `json:"-"`
	Like   int8  `json:"like" validate:"required,oneof=1 2"`
}
