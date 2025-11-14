package core

import (
	"context"

	"github.com/ayayaakasvin/auth-service/internal/models"
)

type UserRepository interface {
	GetPublicUserInfo(ctx context.Context, userID uint) 	(*models.User, error)
	GetPrivateUserInfo(ctx context.Context, userID uint) 	(*models.User, error)
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, username, hashedPassword string) error
	AuthentificateUser(ctx context.Context, username, password string) (uint, error)
}

type Repository interface {
	Close() error

	UserRepository // CRUD

	AuthRepository // Session || Auth
}
