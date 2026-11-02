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
	pool, err := pgxpoohttps://github.com/Mrcoderv/Go_Commerce/pull/6/conflict?name=internal%252Finfrastructure%252Fdatabase%252Fpostgres.go&ancestor_oid=1fab7c42260d8474aadd3857226c9f4cedc4a24b&base_oid=beb717fa59a799daf1c98bfd3a6fc6a840aa5a72&head_oid=c6732bec2e416bbe25a950d326d2cc7065ec6b22l.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Postgres{DB: pool}, nil
}
	
}
