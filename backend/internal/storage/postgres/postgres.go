package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*PGStorage, error) {
	const path = "repository.postgres.New"
	
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	
	return &PGStorage{
		pool: pool,
	}, nil
}

func Close(ctx context.Context, storage *PGStorage) {
	storage.pool.Close()
}