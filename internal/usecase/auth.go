package usecase

import (
	"context"
	"harmoni/internal/conf"
	"harmoni/internal/entity"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/common"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	claimsPool = sync.Pool{
		New: func() any {
			return &entity.JwtCustomClaims{}
		},
	}
)

func resetClaims(claims *entity.JwtCustomClaims) {
	claims.TokenInfo = entity.TokenInfo{}
	claims.RegisteredClaims = jwt.RegisteredClaims{}
}

type AuthUseCase struct {
	userRepo userentity.UserRepository
	logger   *zap.SugaredLogger
	conf     *conf.Auth
}

func NewAuthUseCase(conf *conf.Auth, userRepo userentity.UserRepository, logger *zap.SugaredLogger) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		logger:   logger,
		conf:     conf,
	}
}

func newJwtCustomClaims(userID int64, name string, expiredAt time.Time) *entity.JwtCustomClaims {
	claims := claimsPool.Get().(*entity.JwtCustomClaims)

	claims.TokenInfo = entity.TokenInfo{
		UserID: userID,
		Name:   name,
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	return claims
}

func (u *AuthUseCase) GenToken(ctx context.Context, userID int64, name string) (string, error) {
	claims := newJwtCustomClaims(userID, name, time.Now().Add(time.Duration(u.conf.TokenExpire)*time.Second))
	defer func() {
		resetClaims(claims)
		claimsPool.Put(claims)
	}()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(common.StringToBytes(u.conf.Secret))
	if err != nil {
		return "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tokenStr, nil
}

func (u *AuthUseCase) VerifyToken(ctx context.Context, token string) (jwt.Claims, error) {
	userClaims := claimsPool.Get().(*entity.JwtCustomClaims)
	defer func() {
		resetClaims(userClaims)
		claimsPool.Put(userClaims)
	}()

	tokened, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return common.StringToBytes(u.conf.Secret), nil
	})
	if err != nil {
		return nil, errorx.Unauthorized(reason.TokenInvalid)
	}

	if !tokened.Valid {
		return nil, errorx.Unauthorized(reason.UnauthorizedError)
	}

	return &entity.JwtCustomClaims{
		TokenInfo:        userClaims.TokenInfo,
		RegisteredClaims: userClaims.RegisteredClaims,
	}, nil
}

func (u *AuthUseCase) Revoke() {}
