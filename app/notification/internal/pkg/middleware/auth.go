package middleware

import (
	"context"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/notification/internal/pkg/response"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type currentUserKey struct{}

type AuthUserMiddleware struct {
	uc     v1.UserClient
	logger *log.Helper
}

func NewAuthUserMiddleware(
	uc v1.UserClient,
	logger log.Logger,
) *AuthUserMiddleware {
	return &AuthUserMiddleware{
		uc:     uc,
		logger: log.NewHelper(log.With(logger, "module", "middleware/auth")),
	}
}

func (am *AuthUserMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		if len(token) == 0 {
			c.Next()
			return
		}
		ctx := c.Request.Context()
		resp, err := am.uc.VerifyToken(ctx, &v1.TokenRequest{Token: token})
		if err != nil {
			am.logger.Error(err)
			c.Next()
			return
		}
		if resp != nil && resp.User != nil {
			ctx = context.WithValue(ctx, currentUserKey{}, resp.User)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func (am *AuthUserMiddleware) MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		if len(token) == 0 {
			response.HandleResponse(c, errorx.Unauthorized(reason.UnauthorizedError), nil)
			c.Abort()
			return
		}
		ctx := c.Request.Context()
		resp, err := am.uc.VerifyToken(ctx, &v1.TokenRequest{Token: token})
		if err != nil {
			am.logger.Error(err)
			response.HandleResponse(c, errorx.Unauthorized(reason.UnauthorizedError), nil)
			c.Abort()
			return
		}
		if resp != nil && resp.User != nil {
			ctx = context.WithValue(ctx, currentUserKey{}, resp.User)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func ExtractToken(ctx *gin.Context) (token string) {
	token = ctx.GetHeader("Authorization")
	if len(token) == 0 {
		token = ctx.Query("Authorization")
	}
	return strings.TrimPrefix(token, "Bearer ")
}

func GetUserInfoFromContext(c *gin.Context) *v1.UserBasic {
	userInfo := c.Request.Context().Value(currentUserKey{})
	if userInfo == nil {
		return nil
	}
	u, ok := userInfo.(*v1.UserBasic)
	if !ok {
		return nil
	}
	return u
}
