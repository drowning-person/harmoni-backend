package repository

import (
	authentity "harmoni/app/harmoni/internal/entity/auth"
	commententity "harmoni/app/harmoni/internal/entity/comment"
	eamilentity "harmoni/app/harmoni/internal/entity/email"
	fileentity "harmoni/app/harmoni/internal/entity/file"
	followentity "harmoni/app/harmoni/internal/entity/follow"
	likeentity "harmoni/app/harmoni/internal/entity/like"
	postentity "harmoni/app/harmoni/internal/entity/post"
	tagentity "harmoni/app/harmoni/internal/entity/tag"
	uniqueentity "harmoni/app/harmoni/internal/entity/unique"
	userentity "harmoni/app/harmoni/internal/entity/user"

	"harmoni/app/harmoni/internal/repository/auth"
	"harmoni/app/harmoni/internal/repository/comment"
	"harmoni/app/harmoni/internal/repository/email"
	"harmoni/app/harmoni/internal/repository/file"
	"harmoni/app/harmoni/internal/repository/follow"
	"harmoni/app/harmoni/internal/repository/like"
	"harmoni/app/harmoni/internal/repository/post"
	"harmoni/app/harmoni/internal/repository/tag"
	"harmoni/app/harmoni/internal/repository/unique"
	"harmoni/app/harmoni/internal/repository/user"

	"github.com/google/wire"
)

// ProviderSetRepo is providers.
var ProviderSetRepo = wire.NewSet(

	wire.Bind(new(authentity.AuthRepository), new(*auth.AuthRepo)),
	wire.Bind(new(commententity.CommentRepository), new(*comment.CommentRepo)),
	wire.Bind(new(eamilentity.EmailRepo), new(*email.EmailRepo)),
	wire.Bind(new(followentity.FollowRepository), new(*follow.FollowRepo)),
	wire.Bind(new(postentity.PostRepository), new(*post.PostRepo)),
	wire.Bind(new(tagentity.TagRepository), new(*tag.TagRepo)),
	wire.Bind(new(uniqueentity.UniqueIDRepo), new(*unique.UniqueIDRepo)),
	wire.Bind(new(userentity.UserRepository), new(*user.UserRepo)),
	wire.Bind(new(likeentity.LikeRepository), new(*like.LikeRepo)),

	wire.Bind(new(fileentity.FileRepository), new(*file.FileRepo)),

	auth.NewAuthRepo,
	comment.NewCommentRepo,
	email.NewEmailRepo,
	follow.NewFollowRepo,
	post.NewPostRepo,
	tag.NewTagRepo,
	unique.NewUniqueIDRepo,
	user.NewUserRepo,
	like.NewLikeRepo,
	file.NewFileRepository,
)
