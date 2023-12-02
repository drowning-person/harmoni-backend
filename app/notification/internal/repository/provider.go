package repository

import (
	"harmoni/app/notification/internal/entity/notifyconfig"
	"harmoni/app/notification/internal/entity/remind"
	repoconfig "harmoni/app/notification/internal/repository/notifyconfig"
	reporemind "harmoni/app/notification/internal/repository/remind"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	reporemind.New,
	wire.Bind(new(remind.RemindRepository), new(*reporemind.RemindRepo)),
	repoconfig.NewNotifyConfigRepo,
	wire.Bind(new(notifyconfig.NotifyConfigRepository), new(*repoconfig.NotifyConfigRepo)),
)
