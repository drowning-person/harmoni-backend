package usecase

import (
	"harmoni/app/harmoni/internal/pkg/filesystem"
	"harmoni/app/harmoni/internal/usecase/comment"
	"harmoni/app/harmoni/internal/usecase/email"
	"harmoni/app/harmoni/internal/usecase/file"
	"harmoni/app/harmoni/internal/usecase/follow"
	"harmoni/app/harmoni/internal/usecase/like"
	"harmoni/app/harmoni/internal/usecase/post"
	"harmoni/app/harmoni/internal/usecase/tag"
	"harmoni/app/harmoni/internal/usecase/timeline"
	"harmoni/app/harmoni/internal/usecase/user"

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
