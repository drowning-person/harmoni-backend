package like

import "harmoni/internal/types/object"

type LikeCount struct {
	Counts     int64
	ObjectID   int64
	ObjectType object.ObjectType
}
