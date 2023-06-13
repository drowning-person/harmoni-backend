package handler

import "github.com/google/wire"

// ProviderSetHandler is providers.
var ProviderSetHandler = wire.NewSet(
	NewAccountHandler,
	NewCommentHandler,
	NewFollowHandler,
	NewPostHandler,
	NewTagHandler,
	NewUserHandler,
	NewLikeHandler,
)
