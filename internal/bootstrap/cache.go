package bootstrap

import (
	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/ayayaakasvin/auth-service/internal/repository/valkey"
)

func InitCache(cfg *config.Config) (core.Cache, error)  {
	return valkey.NewValkeyClient(cfg.ValkeyConfig)
}