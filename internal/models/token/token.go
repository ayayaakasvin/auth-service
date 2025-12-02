package token

import (
	"time"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID    uint   `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAccessTokenClaims(uID uint, sID string, ttl time.Duration) *AccessTokenClaims {
	return &AccessTokenClaims{
		UserID:           uID,
		SessionID:        sID,
		RegisteredClaims: genericjwtservice.StdClaims(ttl),
	}
}

func NewRefreshTokenClaims(uID uint, ttl time.Duration) *RefreshTokenClaims {
	return &RefreshTokenClaims{
		UserID: uID,
		RegisteredClaims: genericjwtservice.StdClaims(ttl),
	}
}