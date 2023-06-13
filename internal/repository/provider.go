package repository

import (
	datamysql "harmoni/internal/data/mysql"
	dataredis "harmoni/internal/data/redis"
	authentity "harmoni/internal/entity/auth"
	commententity "harmoni/internal/entity/comment"
	eamilentity "harmoni/internal/entity/email"
	followentity "harmoni/internal/entity/follow"
	likeentity "harmoni/internal/entity/like"
	postentity "harmoni/internal/entity/post"
	tagentity "harmoni/internal/entity/tag"
	uniqueentity "harmoni/internal/entity/unique"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/repository/auth"
	"harmoni/internal/repository/comment"
	"harmoni/internal/repository/email"
	"harmoni/internal/repository/follow"
	"harmoni/internal/repository/like"
	"harmoni/internal/repository/post"
	"harmoni/internal/repository/tag"
	"harmoni/internal/repository/unique"
	"harmoni/internal/repository/user"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
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
	wire.Bind(new(likeentity.LikeRepository), new(*like.LikeRepo)),
	wire.Bind(new(redis.UniversalClient), new(*redis.Client)),

	auth.NewAuthRepo,
	comment.NewCommentRepo,
	email.NewEmailRepo,
	follow.NewFollowRepo,
	post.NewPostRepo,
	tag.NewTagRepo,
	unique.NewUniqueIDRepo,
	user.NewUserRepo,
	like.NewLikeRepo,
)
