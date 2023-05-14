package service

import "github.com/google/wire"

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	NewAccountService,
	NewCommentService,
	NewFollowService,
	NewPostService,
	NewTagService,
	NewUserService,
)
