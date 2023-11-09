package entity

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	TokenInfo
	jwt.RegisteredClaims
}

type TokenInfo struct {
	UserID int64
	Name   string
}
