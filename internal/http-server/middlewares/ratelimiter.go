package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/models/response"
)

var (
	expTimeForRateLimit time.Duration = time.Second * 3
)

const (
	ratelimitformatstring = "ratelimit:%s/%s" // where %d is ip -> plan is like SET: ratelimitformatstring -> true
	loginAction           = "login"
	registerAction        = "login"
)

func (m *Middlewares) RateLimitLoginMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		key := fmt.Sprintf(ratelimitformatstring, loginAction, ip)
		set, err := m.cache.SetNX(r.Context(), key, "1", expTimeForRateLimit)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		if !set {
			response.SendErrorJson(w, http.StatusTooManyRequests, "rate limit, try again in %s seconds", expTimeForRateLimit.String())
			return
		}

		h.ServeHTTP(w, r)
	}
}

func (m *Middlewares) RateLimitRegisterMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		key := fmt.Sprintf(ratelimitformatstring, registerAction, ip)
		set, err := m.cache.SetNX(r.Context(), key, "1", expTimeForRateLimit)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		if !set {
			response.SendErrorJson(w, http.StatusTooManyRequests, "rate limit, try again in %s seconds", expTimeForRateLimit.String())
			return
		}

		h.ServeHTTP(w, r)
	}
}
