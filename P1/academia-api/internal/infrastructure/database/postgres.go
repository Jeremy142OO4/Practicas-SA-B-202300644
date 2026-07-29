package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/config"
)

func NewPostgresPool(
	ctx context.Context,
	cfg config.DatabaseConfig,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(
		cfg.URL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"analizar DATABASE_URL: %w",
			err,
		)
	}

	poolConfig.MaxConns = cfg.MaxConnections
	poolConfig.MinConns = cfg.MinConnections
	poolConfig.ConnConfig.ConnectTimeout =
		cfg.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(
		ctx,
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"crear pool de PostgreSQL: %w",
			err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"comprobar conexión con PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
