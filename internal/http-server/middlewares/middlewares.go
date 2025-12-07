package middlewares

import (
	"strings"

	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/ayayaakasvin/auth-service/internal/services/jwtservice"
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache      core.Cache
	logger     *logrus.Logger
	jwtManager *jwtservice.JWTService

	allowedOrigins   string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool
}

func NewHTTPMiddlewares(logger *logrus.Logger, corsCfg config.CorsConfig, cache core.Cache, jwtManager *jwtservice.JWTService) *Middlewares {
	return &Middlewares{
		logger:     logger,
		cache:      cache,
		jwtManager: jwtManager,

		allowedOrigins:   strings.Join(corsCfg.AllowedOrigins, ","),
		allowedMethods:   strings.Join(corsCfg.AllowedMethods, ","),
		allowedHeaders:   strings.Join(corsCfg.AllowedHeaders, ","),
		allowCredentials: corsCfg.AllowedCredentials,
	}
}
