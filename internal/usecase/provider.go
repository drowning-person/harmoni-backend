package usecase

import (
	"harmoni/internal/pkg/filesystem"
	"harmoni/internal/usecase/file"

	"github.com/google/wire"
)

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
	file.NewPolicy,
	filesystem.NewFileSystem,
	file.NewFileUseCase,
)
