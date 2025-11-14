package postgresql

import (
	"context"

	"github.com/ayayaakasvin/auth-service/internal/models"
)

func (p *PostgreSQL) GetPublicUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	var userObj *models.User = new(models.User)
	userObj.ID = userID
	err := p.conn.QueryRowContext(ctx, "SELECT FROM users (username, created_at) WHERE user_id = $1", userID).Scan(
		&userObj.Username, &userObj.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return userObj, nil
}

func (p *PostgreSQL) GetPrivateUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	var userObj *models.User = new(models.User)
	userObj.ID = userID
	err := p.conn.QueryRowContext(ctx, "SELECT FROM users (username, password, created_at) WHERE user_id = $1", userID).Scan(
		&userObj.Username, &userObj.PasswordHash, &userObj.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return userObj, nil
}
