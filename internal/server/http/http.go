package http

import (
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/iface"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSetHTTP = wire.NewSet(
	NewHTTPServer,
	NewHarmoniAPIRouter,
)

var _ iface.Executor = (*FiberExecutor)(nil)

type FiberExecutor struct {
	conf *config.App
	*fiber.App
}

func (r *FiberExecutor) Start() error {
	return r.Listen(r.conf.Addr)
}

// NewHTTPServer new http server.
func NewHTTPServer(
	conf *config.App,
	zapLogger *zap.Logger,
	harmoniRouter *HarmoniAPIRouter,
	authMiddleware *middleware.JwtAuthMiddleware,
) *FiberExecutor {
	r := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			switch e := err.(type) {
			case *fiber.Error:
				err = errorx.New(e.Code, reason.ServerError).WithMsg(e.Message)
			}
			return fiberx.HandleResponse(c, err, nil)
		},
	})

	r.Use(cors.New())
	recoverconf := recover.ConfigDefault
	recoverconf.EnableStackTrace = true
	r.Use(recover.New(recoverconf))
	r.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 1
	}))

	r.Use(fiberzap.New(fiberzap.Config{
		Logger: zapLogger,
		Fields: []string{"latency", "status", "method", "url", "ip", "call"},
	}))

	harmoniRouter.RegisterStaticRouter(r.Group(""))

	unauthV1 := r.Group("/api/v1")
	unauthV1.Use(authMiddleware.Auth())
	harmoniRouter.RegisterUnAuthHarmoniAPIRouter(unauthV1)

	authV1 := r.Group("/api/v1")
	authV1.Use(authMiddleware.MustAuth())
	harmoniRouter.RegisterHarmoniAPIRouter(authV1)

	return &FiberExecutor{
		conf: conf,
		App:  r,
	}
}
