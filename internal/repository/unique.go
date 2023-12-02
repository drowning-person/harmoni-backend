package repository

import (
	"context"
	"harmoni/internal/types/iface"

	"github.com/bwmarrin/snowflake"
)

var _ iface.UniqueIDRepository = (*UniqueIDRepo)(nil)

type UniqueIDRepo struct {
	node *snowflake.Node
}

func NewUniqueIDRepo(node *snowflake.Node) *UniqueIDRepo {
	return &UniqueIDRepo{node: node}
}

func (r *UniqueIDRepo) GenUniqueID(ctx context.Context) (int64, error) {
	return r.node.Generate().Int64(), nil
}
