package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	pool *pgxpool.Pool
	mx *sync.Mutex
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.Connect(ctx, storagePath)
	if err != nil {
		logrus.Error("failed to create storage")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		pool: pool,
		mx: &sync.Mutex{},
	}, nil
}

func Close(ctx context.Context, storage *Storage) {
	storage.pool.Close()
}