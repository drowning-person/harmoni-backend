package server

import (
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/router"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// NewHTTPServer new http server.
func NewHTTPServer(debug bool,
	zapLogger *zap.Logger,
	harmoniRouter *router.HarmoniAPIRouter,
	authMiddleware *middleware.JwtAuthMiddleware,
) *fiber.App {
	r := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiberx.HandleResponse(c, err, nil)
		},
	})

	r.Use(cors.New())
	r.Use(recover.New())
	r.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 1
	}))

	r.Use(fiberzap.New(fiberzap.Config{
		Logger: zapLogger,
		Fields: []string{"latency", "status", "method", "url", "ip", "call"},
	}))

	unauthV1 := r.Group("/api/v1")
	harmoniRouter.RegisterUnAuthHarmoniAPIRouter(unauthV1)

	authV1 := r.Group("/api/v1")
	authV1.Use(authMiddleware.Auth())
	harmoniRouter.RegisterHarmoniAPIRouter(authV1)

	return r
}
