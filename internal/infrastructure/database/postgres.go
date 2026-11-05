package database

import (
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/migrations"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
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

func NewSQLDB(p *Postgres) *sql.DB {
	return stdlib.OpenDBFromPool(p.DB)
}

func NewFx(lc fx.Lifecycle, cfg *config.Config) (*Postgres, *sql.DB, error) {
	db, err := New(context.Background(), cfg)
	if err != nil {
		return nil, nil, err
	}
	sqlDB := NewSQLDB(db)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.Ping(ctx); err != nil {
				db.Close()
				return err
			}
			if err := migrations.Run(ctx, sqlDB); err != nil {
				db.Close()
				return fmt.Errorf("run migrations: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = sqlDB.Close()
			db.Close()
			return nil
		},
	})

	return db, sqlDB, nil
}

var Module = fx.Module("database", fx.Provide(NewFx))
