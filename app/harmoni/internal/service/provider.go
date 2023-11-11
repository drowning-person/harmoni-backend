package service

import (
	"harmoni/app/harmoni/internal/service/user"

	"github.com/google/wire"
)

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	NewAccountService,
	NewCommentService,
	NewFollowService,
	NewPostService,
	NewTagService,
	NewLikeUsecase,
	NewTimeLineService,
	NewFileService,
	user.ProviderSetUserService,
)
