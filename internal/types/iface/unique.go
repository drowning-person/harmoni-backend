package iface

import "context"

// UniqueIDRepo unique id repository
type UniqueIDRepository interface {
	GenUniqueID(ctx context.Context) (uniqueID int64, err error)
}
