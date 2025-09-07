package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a new pgx connection pool using the provided DSN.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
