package router

import (
	"harmoni/apis"

	"harmoni/pkg/zap"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
)

func New() *fiber.App {
	auth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Error": "请先登录"})
		},
		SigningKey: []byte("secret"),
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{"code": e.Code, "msg": e.Message})
			}
			zap.Logger.Error(err.Error())
			return err
		},
	})
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.Logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "call"},
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 1
	}))
	app.Post("/regist", apis.Regist)
	app.Post("/login", apis.Login)

	apiv1 := app.Group("/apis/v1")

	{
		appUser := apiv1.Group("/user")
		appUser.Get(":id", apis.GetUser)

		appUserWithAuth := appUser.Use(auth)
		appUserWithAuth.Get("", apis.GetAllUsers)
	}
	{
		appTag := apiv1.Group("/tag")
		appTag.Get("", apis.GetTags)
		appTag.Get(":id", apis.GetTagDetail)

		appTagWithAuth := appTag.Use(auth)
		appTagWithAuth.Post("", apis.CreateTag)
	}
	{
		appPost := apiv1.Group("/post")
		appPost.Get(":id", apis.GetPostDetail)
		appPost.Get("", apis.GetPosts)

		appPostWithAuth := appPost.Use(auth)
		appPostWithAuth.Post("", apis.CreatePost)
		appPostWithAuth.Post("/like", apis.LikePost)
	}
	{
		appComment := apiv1.Group("/comment")
		appComment.Get("", apis.GetPostComment)

		appCommentWithAuth := appComment.Use(auth)
		appCommentWithAuth.Post("", apis.CreateComment)
	}
	return app
}
