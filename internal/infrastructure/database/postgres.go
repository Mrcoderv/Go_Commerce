package database

import (
	"Go_Commerce/internal/infrastructure/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
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

func (p *Postgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.DB.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
func NewFx(lc fx.Lifecycle, cfg *config.Config) (*Postgres, error) {
	db, err := New(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.Ping(ctx); err != nil {
				db.Close()
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			db.Close()
			return nil
		},
	})

	return db, nil
}

var module = fx.Module("database", fx.Provide(NewFx))
