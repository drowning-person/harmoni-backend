package repository

import (
	"harmoni/internal/types/iface"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUniqueIDRepo,
	wire.Bind(new(iface.UniqueIDRepository), new(*UniqueIDRepo)),
)
