package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayayaakasvin/auth-service/internal/libs/bcrypt"
	"github.com/ayayaakasvin/auth-service/internal/models"
)

func (p *PostgreSQL) RegisterUser(ctx context.Context, username, hashedPassword string) error {
	if exists, err := p.UsernameExists(ctx, username); err != nil {
		return err
	} else if exists {
		return errors.New("username already in use")
	}

	_, err := p.conn.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) AuthentificateUser(ctx context.Context, username, password string) (uint, error) {
	var (
		userId         uint
		role           models.Role
		hashedPassword string
	)

	if err := p.conn.QueryRowContext(ctx, "SELECT user_id, role, password FROM users WHERE username = $1", username).Scan(&userId, &role, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(NotFound)
		}
		return 0, err
	}

	if err := bcrypt.ComparePasswordAndHash(password, hashedPassword); err != nil {
		return 0, errors.New(UnAuthorized)
	}

	return userId, nil
}

func (p *PostgreSQL) UsernameExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := p.conn.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)", name).Scan(&exists)

	return exists, err
}
