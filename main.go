package main

import (
	"encoding/json"
	"fiberLearn/apis"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		data, _ := json.Marshal(map[string]interface{}{
			"string": "Hello, World!",
			"number": rand.Intn(100),
			"date":   time.Now().Format(time.RFC3339),
		})
		return c.SendString(string(data))
	})
	app.Post("/regist", apis.Regist)
	app.Post("/login", apis.Login)
	{
		appUser := app.Group("/apis/v1")
		appUser.Use(jwtware.New(jwtware.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Error": "请先登录"})
			},
			SigningKey: []byte("secret"),
		}))
		appUser.Get("/user", apis.GetAllUsers)
		appUser.Get("/user/:id", apis.GetUser)
	}
	app.Listen(":80")
}
