package auth

import (
	"context"
	"fmt"
	authentity "harmoni/app/harmoni/internal/entity/auth"
	"harmoni/app/harmoni/internal/pkg/errorx"
	"harmoni/app/harmoni/internal/pkg/reason"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	tokenCanStoreCount = 5
	userPrefix         = "user:"
)

var _ authentity.AuthRepository = (*AuthRepo)(nil)

type AuthRepo struct {
	rdb    *redis.Client
	logger *zap.SugaredLogger
}

func accessTokenKey(userID int64, token string) string {
	return fmt.Sprintf("%s%d:token:%s", userPrefix, userID, token)
}

func refreshTokenKey(userID int64) string {
	return fmt.Sprintf("%s%d:refresh_token", userPrefix, userID)
}

func NewAuthRepo(rdb *redis.Client, logger *zap.SugaredLogger) *AuthRepo {
	return &AuthRepo{
		rdb:    rdb,
		logger: logger.With("module", "repository/auth"),
	}
}

// storeAccessToken does not need to store access token with JWT for now
func (r *AuthRepo) storeAccessToken(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := accessTokenKey(userID, token)
	_, err := r.rdb.Set(ctx, key, "", ttl).Result()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *AuthRepo) storeRefreshToken(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := refreshTokenKey(userID)
	tokenCount, err := r.rdb.ZCard(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errorx.BadRequest(reason.TokenNotExistOrExpired)
		}
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	expireAt := time.Now().Add(ttl).Unix()
	if tokenCount > 5 {
		_, err := r.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
			_, err := p.ZPopMin(ctx, key, 1).Result()
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}

	member := redis.Z{
		Score:  float64(expireAt), // 最后使用时间
		Member: token,
	}
	count, err := r.rdb.ZAddNX(ctx, key, member).Result()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if count != 1 {
		r.logger.Warnf("store token failed,count should be one")
	}

	return nil
}

// StoreToken store a token ID
// The "token" parameter in this function is actually the JWT token ID.
func (r *AuthRepo) StoreToken(ctx context.Context, userID int64, token string, tokenType authentity.TokenType, ttl time.Duration) error {
	switch tokenType {
	case authentity.AccessTokenType:
		return r.storeAccessToken(ctx, userID, token, ttl)
	case authentity.RefreshTokenType:
		return r.storeRefreshToken(ctx, userID, token, ttl)
	}

	return nil
}

// RetrieveToken retrieve token from redis
// The "token" parameter in this function is actually the JWT token ID.
func (r *AuthRepo) RetrieveToken(ctx context.Context, userID int64, token string, tokenType authentity.TokenType) (string, bool, error) {
	var (
		key string
	)
	switch tokenType {
	case authentity.AccessTokenType:
		key = accessTokenKey(userID, token)
		if count, err := r.rdb.Exists(ctx, key).Result(); err != nil {
			return "", false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if count == 1 {
			return "", true, nil
		}
	case authentity.RefreshTokenType:
		key = refreshTokenKey(userID)
		_, err := r.rdb.ZScore(ctx, key, token).Result()
		if err != nil {
			if err == redis.Nil {
				return "", false, nil
			}
			return "", false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return "", true, nil
	}

	return "", false, nil
}

// RetrieveTokens do nothing for now
func (r *AuthRepo) RetrieveTokens(ctx context.Context, userID int64, tokenType authentity.TokenType) ([]string, bool, error) {
	return nil, false, nil
}

// RemoveToken removes a token and adds it to the blacklist.
// The "token" parameter in this function is actually the JWT token ID.
// The "ttl" parameter in this function is the token on black list ttl
func (r *AuthRepo) RemoveToken(ctx context.Context, userID int64, token string, tokenType authentity.TokenType, ttl time.Duration) error {
	var (
		count int64
		err   error
	)
	switch tokenType {
	case authentity.AccessTokenType:
		key := accessTokenKey(userID, token)
		count, err = r.rdb.Unlink(ctx, key).Result()
	case authentity.RefreshTokenType:
		key := refreshTokenKey(userID)
		count, err = r.rdb.ZRem(ctx, key, token).Result()
	}

	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if count != 1 {
		r.logger.Warnf("remove token failed,count should be one.Token type:%#v UserID: %d", tokenType, userID)
		return nil
	}

	return nil
}

func (r *AuthRepo) removeAccessTokens(ctx context.Context, userID int64, ttl time.Duration) error {
	key := fmt.Sprintf("%s%d:*", userPrefix, userID)
	for {
		tokens, cursor, err := r.rdb.Scan(ctx, 0, key, 10).Result()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		count, err := r.rdb.Unlink(ctx, tokens...).Result()
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		} else if count != int64(len(key)) {
			r.logger.Warnf("remove access token failed, count is not equal to len of keys.count: %v length: %v", count, len(key))
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *AuthRepo) removeRefreshTokens(ctx context.Context, userID int64, ttl time.Duration) error {
	key := refreshTokenKey(userID)
	tokens, cursor, err := r.rdb.ZScan(ctx, key, 0, "", tokenCanStoreCount).Result()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if cursor != 0 {
		r.logger.Warnf("scan refresh token zset failed, cursor is not zero but %d", cursor)
	}

	_, err = r.rdb.TxPipelined(ctx, func(p redis.Pipeliner) error {
		count, err := p.ZRem(ctx, key, tokens).Result()
		if err != nil {
			return err
		} else if count != int64(len(tokens)) {
			r.logger.Warnf("remove token count is not equal to zscan but %d, zscan count: %d", count, len(tokens))
		}

		return nil
	})

	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *AuthRepo) RemoveTokens(ctx context.Context, userID int64, tokenType authentity.TokenType, ttl time.Duration) error {
	switch tokenType {
	case authentity.AccessTokenType:
		return r.removeAccessTokens(ctx, userID, ttl)
	case authentity.RefreshTokenType:
		return r.removeRefreshTokens(ctx, userID, ttl)
	}

	return nil
}
