package unique

import (
	"context"
	"harmoni/app/harmoni/internal/entity/unique"

	"github.com/bwmarrin/snowflake"
)

var _ unique.UniqueIDRepo = (*UniqueIDRepo)(nil)

type UniqueIDRepo struct {
	node *snowflake.Node
}

func NewUniqueIDRepo(node *snowflake.Node) *UniqueIDRepo {
	return &UniqueIDRepo{node: node}
}

func (r *UniqueIDRepo) GenUniqueID(ctx context.Context) (int64, error) {
	return r.node.Generate().Int64(), nil
}
