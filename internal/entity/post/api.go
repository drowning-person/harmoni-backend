package post

import (
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
)

type GetPostsRequest struct {
	entity.PageCond
	UserID int64  `json:"-"`
	Order  string `query:"order" validate:"omitempty,oneof=like rtime ctime" label:"排序"`
	TagID  int64  `query:"tag_id"`
}

type GetPostsReply struct {
	*paginator.Page[PostBasicInfo] `json:"posts"`
}

type GetPostInfoRequest struct {
	UserID int64 `json:"-"`
	PostID int64 `params:"id" label:"帖子ID"`
}

type GetPostInfoReply struct {
	PostInfo
}

type CreatePostRequest struct {
	TagIDs  entity.Int64Slice `json:"tag_ids" validate:"lte=4" label:"话题ID"`
	UserID  int64             `json:"-"`
	Title   string            `json:"title" validate:"required,gte=3,lte=128" label:"帖子标题"`
	Content string            `json:"content" validate:"required,gte=6,lte=65535" label:"帖子内容"`
}

type CreatePostReply struct {
	PostInfo
}

type LikePostRequest struct {
	PostID int64 `json:"post_id,string" validate:"required" label:"帖子ID"`
	UserID int64 `json:"-"`
	Like   int8  `json:"like" validate:"required,oneof=1 2"`
}
