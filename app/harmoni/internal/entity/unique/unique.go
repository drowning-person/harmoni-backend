package unique

import "context"

// UniqueIDRepo unique id repository
type UniqueIDRepo interface {
	GenUniqueID(ctx context.Context) (uniqueID int64, err error)
}
