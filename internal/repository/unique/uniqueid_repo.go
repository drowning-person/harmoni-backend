package unique

import (
	"context"
	"harmoni/internal/entity/unique"

	"github.com/bwmarrin/snowflake"
)

type uniqueIDRepo struct {
	node *snowflake.Node
}

func NewUniqueIDRepo(node *snowflake.Node) unique.UniqueIDRepo {
	return &uniqueIDRepo{node: node}
}

func (r *uniqueIDRepo) GenUniqueID(ctx context.Context) (int64, error) {
	return r.node.Generate().Int64(), nil
}
