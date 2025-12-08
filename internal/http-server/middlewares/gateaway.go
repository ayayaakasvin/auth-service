package middlewares

import (
	"net/http"

	"github.com/ayayaakasvin/auth-service/internal/models/response"
)

const gateAwayHeader = "X-Gateaway-Key"

// Since most of mine services are hosted in Render free instance, GateAwayMiddleware is used to make service reachable only via gateaway
func (mw *Middlewares) GateAwayMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(gateAwayHeader) != "" {
			response.SendErrorJson(w, http.StatusForbidden, "Forbidden")
			return 
		}

		next(w, r)
	}
}