package repository

import (
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/app/like/internal/repository/like"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	like.NewLikeRepo,
	wire.Bind(new(entitylike.LikeRepository), new(*like.LikeRepo)),
)
