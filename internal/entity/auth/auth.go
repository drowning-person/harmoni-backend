package auth

import (
	"context"
	"time"
)

type TokenType uint8

const (
	AccessTokenType TokenType = iota + 1
	RefreshTokenType
)

type AuthRepository interface {
	StoreToken(ctx context.Context, userID int64, token string, tokenType TokenType, ttl time.Duration) error
	RetrieveToken(ctx context.Context, userID int64, token string, tokenType TokenType) (string, bool, error)
	RetrieveTokens(ctx context.Context, userID int64, tokenType TokenType) ([]string, bool, error)
	RemoveToken(ctx context.Context, userID int64, token string, tokenType TokenType, ttl time.Duration) error
	RemoveTokens(ctx context.Context, userID int64, tokenType TokenType, ttl time.Duration) error
}
