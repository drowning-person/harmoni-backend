package repository

import (
	datamysql "harmoni/internal/data/mysql"
	dataredis "harmoni/internal/data/redis"
	authentity "harmoni/internal/entity/auth"
	commententity "harmoni/internal/entity/comment"
	eamilentity "harmoni/internal/entity/email"
	followentity "harmoni/internal/entity/follow"
	postentity "harmoni/internal/entity/post"
	tagentity "harmoni/internal/entity/tag"
	uniqueentity "harmoni/internal/entity/unique"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/repository/auth"
	"harmoni/internal/repository/comment"
	"harmoni/internal/repository/email"
	"harmoni/internal/repository/follow"
	"harmoni/internal/repository/post"
	"harmoni/internal/repository/tag"
	"harmoni/internal/repository/unique"
	"harmoni/internal/repository/user"

	"github.com/google/wire"
)

// ProviderSetRepo is providers.
var ProviderSetRepo = wire.NewSet(
	dataredis.NewRedis,
	datamysql.NewDB,
	wire.Bind(new(authentity.AuthRepository), new(*auth.AuthRepo)),
	wire.Bind(new(commententity.CommentRepository), new(*comment.CommentRepo)),
	wire.Bind(new(eamilentity.EmailRepo), new(*email.EmailRepo)),
	wire.Bind(new(followentity.FollowRepository), new(*follow.FollowRepo)),
	wire.Bind(new(postentity.PostRepository), new(*post.PostRepo)),
	wire.Bind(new(tagentity.TagRepository), new(*tag.TagRepo)),
	wire.Bind(new(uniqueentity.UniqueIDRepo), new(*unique.UniqueIDRepo)),
	wire.Bind(new(userentity.UserRepository), new(*user.UserRepo)),

	auth.NewAuthRepo,
	comment.NewCommentRepo,
	email.NewEmailRepo,
	follow.NewFollowRepo,
	post.NewPostRepo,
	tag.NewTagRepo,
	unique.NewUniqueIDRepo,
	user.NewUserRepo,
)
