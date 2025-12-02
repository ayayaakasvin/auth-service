package jwtservice

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/config"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	*genericjwtservice.JWTService
	Secret []byte

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTManager(cfg *config.JWTSecret) *JWTService {
	return &JWTService{
		JWTService: genericjwtservice.NewJWTService([]byte(cfg.Secret), jwt.SigningMethodHS256),

		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (j *JWTService) FetchUserID(userIdAny any) (uint, error) {
	switch v := userIdAny.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case string:
		idInt, err := strconv.Atoi(v)
		return uint(idInt), err
	default:
		return 0, fmt.Errorf("invalid user id type")
	}
}