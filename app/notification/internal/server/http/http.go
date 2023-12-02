package http

import (
	"harmoni/app/notification/internal/pkg/middleware"
	"harmoni/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/tx7do/kratos-transport/transport/gin"
)

var ProviderSetHTTP = wire.NewSet(
	NewNotificationAPIRouter,
	NewHTTPServer,
)

func NewHTTPServer(
	cong *conf.Server,
	logger log.Logger,
	authMiddleware *middleware.AuthUserMiddleware,
	router *NotificationAPIRouter,
) *gin.Server {
	server := gin.NewServer(gin.WithAddress(cong.Http.Addr))
	server.Use(gin.GinLogger(logger))
	server.Use(gin.GinRecovery(logger, true))
	authGroup := server.Group("/api/v1", authMiddleware.MustAuth())
	router.RegisterNotificationAPIRouter(authGroup)
	unAuthGroup := server.Group("/api/v1", authMiddleware.Auth())
	router.RegisterUnAuthNotificationAPIRouter(unAuthGroup)
	return server
}
