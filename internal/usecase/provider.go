package usecase

import "github.com/google/wire"

// ProviderSetUsecase is providers.
var ProviderSetUsecase = wire.NewSet(
	NewAccountUsecase,
	NewAuthUseCase,
	NewCommentUseCase,
	NewEmailUsecase,
	NewFollowUseCase,
	NewPostUseCase,
	NewTagUseCase,
	NewUserUseCase,
	NewLikeUsecase,
	NewTimeLineUsecase,
)
