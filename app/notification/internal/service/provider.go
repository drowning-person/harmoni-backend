package service

import (
	"harmoni/app/notification/internal/service/remind"

	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(
	remind.NewRemindService,
)
