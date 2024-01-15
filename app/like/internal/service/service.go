package service

import (
	"harmoni/app/like/internal/service/like"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	like.NewLikeService,
)
