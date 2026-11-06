package database

import (
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"
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

func NewFx(lc fx.Lifecycle, cfg *config.Config, log logger.Logger) (*Postgres, *sql.DB, error) {
	log.Info("connecting to postgres", "host", cfg.DBHost, "port", cfg.DBPort, "database", cfg.DBName)
	db, err := New(context.Background(), cfg) //  the new  db context is borned
	if err != nil {
		log.Error("postgres connection failed", "error", err)
		return nil, nil, err
	}
	sqlDB := NewSQLDB(db)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("pinging postgres")
			if err := db.Ping(ctx); err != nil {
				log.Error("postgres ping failed", "error", err)
				db.Close()
				return err
			}
			log.Info("running database migrations")
			if err := migrations.Run(ctx, sqlDB); err != nil {
				log.Error("database migrations failed", "error", err)
				db.Close()
				return fmt.Errorf("run migrations: %w", err)
			}
			log.Info("database ready", "host", cfg.DBHost, "port", cfg.DBPort, "database", cfg.DBName)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("closing database")
			_ = sqlDB.Close()
			db.Close()
			return nil
		},
	})

	return db, sqlDB, nil
}

var Module = fx.Module("database", fx.Provide(NewFx))
