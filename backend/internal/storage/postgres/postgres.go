package postgres

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*PGStorage, error) {
	const path = "repository.postgres.New"

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	return &PGStorage{
		pool: pool,
	}, nil
}

func Close(ctx context.Context, storage *PGStorage) {
	storage.pool.Close()
}
