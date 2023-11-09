package middleware

import (
	"context"
	"harmoni/app/harmoni/internal/entity"
	authentity "harmoni/app/harmoni/internal/entity/auth"
	"harmoni/app/harmoni/internal/pkg/errorx"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/app/harmoni/internal/usecase/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenCtxType string

type Config struct {
	Unauthorized  func(error) error
	TokenCtxName  TokenCtxType
	TokenHeadName string
	TokenLookup   string
}

type JwtAuthMiddleware struct {
	conf        *Config
	authUsecase *user.AuthUseCase
}

var (
	defaultConf = &Config{
		TokenLookup:   "header: Authorization, query: token, cookie: token",
		TokenHeadName: "Bearer",
		TokenCtxName:  "TOKEN_CTX",
	}
)

type Option func(conf *Config)

func NewJwtAuthMiddleware(authUsecase *user.AuthUseCase) *JwtAuthMiddleware {
	mw := &JwtAuthMiddleware{
		conf:        defaultConf,
		authUsecase: authUsecase,
	}

	/* 	for _, opt := range opts {
		opt(mw.conf)
	} */

	return mw
}

func (mw *JwtAuthMiddleware) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := mw.ParseAndVerifyToken(c, false)
		if err != nil {
			return err
		}

		return c.Next()
	}
}

func (mw *JwtAuthMiddleware) MustAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := mw.ParseAndVerifyToken(c, true)
		if err != nil {
			return err
		}

		return c.Next()
	}
}

func (mw *JwtAuthMiddleware) jwtFromHeader(c *fiber.Ctx, key string) (string, error) {
	authHeader := c.Context().Request.Header.Peek(key)

	if string(authHeader) == "" {
		return "", errorx.Unauthorized(reason.UnauthorizedError)
	}

	parts := strings.SplitN(string(authHeader), " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.conf.TokenHeadName) {
		return "", errorx.Unauthorized(reason.AuthHeaderInvalid)
	}

	return parts[1], nil
}

func (mw *JwtAuthMiddleware) jwtFromQuery(c *fiber.Ctx, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		return "", errorx.Unauthorized(reason.UnauthorizedError)
	}

	return token, nil
}

func (mw *JwtAuthMiddleware) jwtFromCookie(c *fiber.Ctx, key string) (string, error) {

	cookie := c.Cookies(key)

	if cookie == "" {
		return "", errorx.Unauthorized(reason.UnauthorizedError)
	}

	return cookie, nil
}

func (mw *JwtAuthMiddleware) jwtFromParam(c *fiber.Ctx, key string) (string, error) {
	token := c.Params(key)
	if token == "" {
		return "", errorx.Unauthorized(reason.UnauthorizedError)
	}
	return token, nil
}

func (mw *JwtAuthMiddleware) ParseAndVerifyToken(c *fiber.Ctx, must bool) error {
	var token string
	var err error

	methods := strings.Split(mw.conf.TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = mw.jwtFromHeader(c, v)
		case "query":
			token, err = mw.jwtFromQuery(c, v)
		case "cookie":
			token, err = mw.jwtFromCookie(c, v)
		case "param":
			token, err = mw.jwtFromParam(c, v)
		}
	}

	if err != nil && must {
		return err
	}

	claims, err := mw.authUsecase.VerifyToken(c.UserContext(), token, authentity.AccessTokenType)
	if err != nil && must {
		return err
	}
	if claims != nil {
		ctx := context.WithValue(c.UserContext(), mw.conf.TokenCtxName, claims)
		c.SetUserContext(ctx)
	}

	return nil
}

func WithTokenCtxName(name string) Option {
	return func(conf *Config) {
		conf.TokenCtxName = TokenCtxType(name)
	}
}

func WithTokenHeadName(name string) Option {
	return func(conf *Config) {
		conf.TokenHeadName = name
	}
}

func WithTokenLookup(method string) Option {
	return func(conf *Config) {
		conf.TokenLookup = method
	}
}

func GetClaimsFromCtx(ctx context.Context) *entity.JwtCustomClaims {
	tmp := ctx.Value(defaultConf.TokenCtxName)
	if tmp == nil {
		return nil
	}

	claims, ok := tmp.(*entity.JwtCustomClaims)
	if !ok {
		return nil
	}

	return claims
}
