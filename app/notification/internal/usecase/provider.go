package usecase

import (
	"harmoni/app/notification/internal/usecase/remind"
	"harmoni/app/notification/internal/usecase/remind/events"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	remind.NewRemindUsecase,
	events.NewLikeEventsHandler,
)
