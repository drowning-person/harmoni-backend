package user

import (
	"context"
	"harmoni/app/harmoni/internal/entity"
	authentity "harmoni/app/harmoni/internal/entity/auth"
	userentity "harmoni/app/harmoni/internal/entity/user"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/pkg/common"
	"harmoni/app/harmoni/internal/pkg/errorx"
	"harmoni/app/harmoni/internal/pkg/reason"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	authRepo authentity.AuthRepository
	logger   *zap.SugaredLogger
	conf     *config.Auth
}

func NewAuthUseCase(conf *config.Auth, authRepo authentity.AuthRepository, logger *zap.SugaredLogger) *AuthUseCase {
	return &AuthUseCase{
		authRepo: authRepo,
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
		ID:        uuid.New().String(),
	}

	return claims
}

func (u *AuthUseCase) GenToken(ctx context.Context, user *userentity.User, tokenType authentity.TokenType) (string, error) {
	var (
		ttl time.Duration
		ext time.Time
	)

	switch tokenType {
	case authentity.AccessTokenType:
		ttl = u.conf.TokenExpire
		ext = time.Now().Add(ttl)
	case authentity.RefreshTokenType:
		ttl = u.conf.RefreshTokenExpire
		ext = time.Now().Add(ttl)
	}

	claims := newJwtCustomClaims(user.UserID, user.Name, ext)
	defer func() {
		resetClaims(claims)
		claimsPool.Put(claims)
	}()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(common.StringToBytes(u.conf.Secret))
	if err != nil {
		return "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// Store the JWT ID in a secure location
	// This will be used to revoke user access to resources during certain processes
	// such as password change, logout, and exit.
	err = u.authRepo.StoreToken(ctx, user.UserID, claims.ID, tokenType, ttl)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ReGenToken regenerate refresh token
func (u *AuthUseCase) ReGenToken(ctx context.Context, userClaims *entity.JwtCustomClaims, tokenType authentity.TokenType) (string, error) {
	switch tokenType {
	case authentity.RefreshTokenType:
		u.logger.Debugf("refresh token expire at %v", userClaims.ExpiresAt.GoString())
		if time.Until(userClaims.ExpiresAt.Time) <= (u.conf.RefreshTokenExpire / 2) {
			newRefreshToken, err := u.GenToken(ctx, &userentity.User{UserID: userClaims.UserID, Name: userClaims.Name}, authentity.RefreshTokenType)
			if err != nil {
				u.logger.Warnf("generate new refresh token failed %s", err.Error())
			} else {
				return newRefreshToken, nil
			}

		}
	}

	return "", nil
}

// VerifyToken verifies the authenticity and validity of a JWT token
// and retrieves the user claims associated with it.
// It takes in the token string, tokenType (refresh or access), and context object as parameters.
// It returns the user claims if the token is valid and not blacklisted,
// otherwise, it returns an error indicating the reason for failure.
func (u *AuthUseCase) VerifyToken(ctx context.Context, token string, tokenType authentity.TokenType) (*entity.JwtCustomClaims, error) {
	userClaims := &entity.JwtCustomClaims{}

	tokened, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return common.StringToBytes(u.conf.Secret), nil
	})
	if err != nil {
		u.logger.Error(err)
		return nil, errorx.Unauthorized(reason.TokenInvalid)
	}

	if !tokened.Valid {
		return nil, errorx.Unauthorized(reason.UnauthorizedError)
	}

	_, exist, err := u.authRepo.RetrieveToken(ctx, userClaims.UserID, userClaims.ID, tokenType)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, errorx.Unauthorized(reason.TokenInvalid)
	}

	return userClaims, nil
}

// RevokeToken revoke access token or refresh token
func (u *AuthUseCase) RevokeToken(ctx context.Context, userID int64, token string, tokenType authentity.TokenType) error {
	return u.authRepo.RemoveToken(ctx, userID, token, tokenType, u.conf.RefreshTokenExpire)
}

// RevokeTokens revoke all access tokens and refresh tokens
func (u *AuthUseCase) RevokeTokens(ctx context.Context, userID int64, tokenType authentity.TokenType) error {
	return u.authRepo.RemoveTokens(ctx, userID, tokenType, u.conf.RefreshTokenExpire)
}
