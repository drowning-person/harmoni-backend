package like

import (
	entitylike "harmoni/app/like/internal/entity/like"
)

type LikeCacheStrategy interface {
	QueryUserLikeListFromDB(query *entitylike.ListLikeObjectQuery) bool
}

type defaultLikeCacheStrategy int

func (s defaultLikeCacheStrategy) QueryUserLikeListFromDB(query *entitylike.ListLikeObjectQuery) bool {
	return query.End() > int64(s)
}
