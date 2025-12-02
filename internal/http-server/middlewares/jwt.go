package middlewares

import (
	"net/http"
	"strings"

	"github.com/ayayaakasvin/auth-service/internal/ctx"
	"github.com/ayayaakasvin/auth-service/internal/models/response"
	"github.com/ayayaakasvin/auth-service/internal/models/token"
	"github.com/redis/go-redis/v9"
)

const (
	AuthorizationHeader = "Authorization"
)

// JWTAuthMiddleware is a middleware for http.HandlerFunc
func (m *Middlewares) JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			unauthorized(w, "authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			unauthorized(w, "authorization header missing")
			return
		}

		cl, err := m.jwtManager.Validate(tokenString, &token.AccessTokenClaims{})
		if err != nil {
			unauthorized(w, "failed to validate jwt")
			return
		}

		fullClaims, ok := cl.(*token.AccessTokenClaims)
		if !ok {
			unauthorized(w, "invalid claims")
			return
		}

		if fullClaims.SessionID == "" {
			unauthorized(w, "session_id missing")
			return
		}

		if _, err := m.cache.Get(r.Context(), fullClaims.SessionID); err == redis.Nil {
			unauthorized(w, "session is expired")
			return
		}

		if fullClaims.UserID == 0 {
			unauthorized(w, "user_id missing or invalid")
			return
		}

		r = ctx.WrapValueIntoRequest(r, ctx.CtxUserIDKey, fullClaims.UserID)
		r = ctx.WrapValueIntoRequest(r, ctx.CtxSessionIDKey, fullClaims.SessionID)

		next(w, r)
	}
}

func unauthorized(w http.ResponseWriter, msg string) {
	response.SendErrorJson(w, http.StatusUnauthorized, "%s", msg)
}
