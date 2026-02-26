package bootstrap

import (
	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/ayayaakasvin/auth-service/internal/repository/postgresql"
)

func InitRepository(cfg *config.Config) (core.Repository, error) {
    return postgresql.NewPostgreSQLConnection(cfg.PostgreSQLConfig)
}