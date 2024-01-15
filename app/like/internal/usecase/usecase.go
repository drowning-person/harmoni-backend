package usecase

import (
	"harmoni/app/like/internal/usecase/like"
	"harmoni/app/like/internal/usecase/like/events"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	events.NewLikeEventsHandler,
	like.NewLikeUsecase,
)
