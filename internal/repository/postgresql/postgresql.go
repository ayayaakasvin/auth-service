package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models/core"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	conn *sql.DB
}

func NewPostgreSQLConnection(dbConfig config.PostgreSQLConfig) (core.Repository, error) {
	psql := new(PostgreSQL)

	connection, err := sql.Open("postgres", dbConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	psql.conn = connection

	if err := psql.conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping to db: %w", err)
	}

	return psql, nil
}

func (p *PostgreSQL) Close() error {
	return p.Close()
}
