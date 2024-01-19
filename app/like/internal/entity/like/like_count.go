package like

import (
	objectv1 "harmoni/api/common/object/v1"
)

type LikeCount struct {
	Count  int64
	Object *objectv1.Object
}

type LikeCountList []*LikeCount

func (l LikeCountList) ToMap() map[int64]int64 {
	m := make(map[int64]int64)
	for _, v := range l {
		m[v.Object.GetId()] = v.Count
	}
	return m
}
