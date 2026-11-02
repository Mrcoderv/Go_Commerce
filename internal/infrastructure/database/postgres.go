package database

import (
	"Go_Commerce/internal/infrastructure/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Postgres{DB: pool}, nil
}
