package postgresql

import (
	"context"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/models"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
)

var testUserInfo *models.User = &models.User{
	ID:           0,
	Username:     "Jane Doe",
	PasswordHash: "hashedPassword",
	Role:         "test",
	CreatedAt:    time.Unix(0, 0),
}

type PostgreSQL_Mock struct{}

// AuthentificateUser implements core.Repository.
func (p *PostgreSQL_Mock) AuthentificateUser(ctx context.Context, username string, password string) (uint, error) {
	return testUserInfo.ID, nil
}

// Close implements core.Repository.
func (p *PostgreSQL_Mock) Close() error {
	return nil
}

// GetPrivateUserInfo implements core.Repository.
func (p *PostgreSQL_Mock) GetPrivateUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	return testUserInfo, nil
}

// GetPublicUserInfo implements core.Repository.
func (p *PostgreSQL_Mock) GetPublicUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	return testUserInfo, nil
}

// RegisterUser implements core.Repository.
func (p *PostgreSQL_Mock) RegisterUser(ctx context.Context, username string, hashedPassword string) error {
	return nil
}

func NewPostgreSQL_Mock() core.Repository {
	return &PostgreSQL_Mock{}
}
