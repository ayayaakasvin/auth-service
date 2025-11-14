// Handlers that serves for main http server, accessed via handlers.Handler struct that contains necessary dependencies
package handlers

import (
	"github.com/ayayaakasvin/auth-service/internal/jwttool"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	repo		core.Repository
	cache 		core.Cache
	jwtM		*jwttool.JWTManager

	logger 		*logrus.Logger
}

func NewHTTPHandlers(repo core.Repository, cache core.Cache, logger *logrus.Logger, jwtM *jwttool.JWTManager) *Handlers {
	return &Handlers{
		repo: repo,
		cache: cache,
		jwtM: jwtM,

		logger: logger,
	}
}