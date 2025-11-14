package jwttool

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret []byte

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTManager(cfg *config.JWTSecret) *JWTManager {
	return &JWTManager{
		Secret: []byte(cfg.Secret),

		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (j *JWTManager) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWTManager) GenerateAccessToken(userId uint, sessionId string, ttl time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTToken{
		UserID:    userId,
		SessionID: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		log.Print(err)
		return ""
	}

	return tokenString
}

func (j *JWTManager) GenerateRefreshToken(userId uint, ttl time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTToken{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		log.Print(err)
		return ""
	}

	return tokenString
}

func (j *JWTManager) FetchUserID(userIdAny any) (uint, error) {
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
