package server

import (
	"harmoni/app/like/internal/server/mq"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewGRPCServer,
	mq.NewPublisher,
	mq.NewLikeGroup,
	mq.NewMQRouter,
)
