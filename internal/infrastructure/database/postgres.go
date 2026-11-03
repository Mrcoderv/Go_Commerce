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
	pool, err := pgxpoohttps://github.com/Mrcoderv/Go_Commerce/pull/6/conflict?name=internal%252Finfrastructure%252Fdatabase%252Fpostgres.go&ancestor_oid=1fab7c42260d8474aadd3857226c9f4cedc4a24b&base_oid=beb717fa59a799daf1c98bfd3a6fc6a840aa5a72&head_oid=c6732bec2e416bbe25a950d326d2cc7065ec6b22l.New(ctx, cfg.DSN())
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

// NewFx is the constructor Fx calls. It ties ping and close to the app lifecycle.
func NewFx(lc fx.Lifecycle, cfg *config.Config) (*Postgres, error) {
	// context.Background(), not the hook's ctx: the pool must outlive startup.
	db, err := New(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.Ping(ctx); err != nil {
				db.Close() // don't leak the pool if startup fails
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

var Module = fx.Module("database", fx.Provide(NewFx))
	
}
