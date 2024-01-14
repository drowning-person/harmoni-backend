package http

import (
	"harmoni/app/notification/internal/handler/remind"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type NotificationAPIRouter struct {
	rh     *remind.RemindHandler
	logger *log.Helper
}

func NewNotificationAPIRouter(
	rh *remind.RemindHandler,
	logger log.Logger,
) *NotificationAPIRouter {
	return &NotificationAPIRouter{
		rh:     rh,
		logger: log.NewHelper(log.With(logger, "module", "server/http")),
	}
}

func (r *NotificationAPIRouter) RegisterNotificationAPIRouter(g *gin.RouterGroup) {
	mg := g.Group("/message")
	{
		mg.GET("/unread", r.rh.UnReadCount)

		mg.GET("/like", r.rh.ListLikeRemind)
		mg.GET("/like/detail", r.rh.ListLikeRemindDetail)

		mg.GET("/reply", r.rh.ListReplyRemind)
		mg.GET("/at", r.rh.ListAtRemind)
	}
}

func (r *NotificationAPIRouter) RegisterUnAuthNotificationAPIRouter(g *gin.RouterGroup) {
}
