package middlewares

import (
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/ayayaakasvin/auth-service/internal/services/jwtservice"
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache      core.Cache
	logger     *logrus.Logger
	jwtManager *jwtservice.JWTService
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache core.Cache, jwtManager *jwtservice.JWTService) *Middlewares {
	return &Middlewares{
		logger:     logger,
		cache:      cache,
		jwtManager: jwtManager,
	}
}
