package middlewares

import (
	"github.com/ayayaakasvin/auth-service/internal/jwttool"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache 		core.Cache
	logger 		*logrus.Logger
	jwtManager	*jwttool.JWTManager
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache core.Cache, jwtManager	*jwttool.JWTManager) *Middlewares {
	return &Middlewares{
		logger: logger,
		cache: cache,
		jwtManager: jwtManager,
	}
}