package router

import (
	"fiberLearn/apis"

	"fiberLearn/pkg/zap"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	jwtware "github.com/gofiber/jwt/v3"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			e := err.(*fiber.Error)
			return c.Status(e.Code).JSON(fiber.Map{"code": e.Code, "msg": e.Message})
		},
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.Logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "call"},
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 1
	}))
	app.Post("/regist", apis.Regist)
	app.Post("/login", apis.Login)

	appWithAuth := app.Group("/apis/v1", jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Error": "请先登录"})
		},
		SigningKey: []byte("secret"),
	}))

	{
		appUser := appWithAuth.Group("/user")
		appUser.Get("", apis.GetAllUsers)
		appUser.Get(":id", apis.GetUser)
	}
	{
		appTag := appWithAuth.Group("/tag")
		appTag.Get("", apis.GetTags)
		appTag.Get(":id", apis.GetTagDetail)
		appTag.Post("", apis.CreateTag)
	}
	return app
}
