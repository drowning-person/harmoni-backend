package usecase

import (
	"harmoni/internal/pkg/filesystem"
	"harmoni/internal/usecase/comment"
	"harmoni/internal/usecase/email"
	"harmoni/internal/usecase/file"
	"harmoni/internal/usecase/follow"
	"harmoni/internal/usecase/like"
	"harmoni/internal/usecase/post"
	"harmoni/internal/usecase/tag"
	"harmoni/internal/usecase/timeline"
	"harmoni/internal/usecase/user"

	"github.com/google/wire"
)

// ProviderSetUsecase is providers.
var ProviderSetUsecase = wire.NewSet(
	user.ProviderSetUser,
	post.ProviderSetPost,
	like.ProviderSetLikeUsecase,
	comment.ProviderSetComment,

	email.NewEmailUsecase,
	follow.NewFollowUseCase,
	tag.NewTagUseCase,
	timeline.NewTimeLineUsecase,
	file.NewPolicy,
	filesystem.NewFileSystem,
	file.NewFileUseCase,
)
