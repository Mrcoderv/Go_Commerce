package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

// using the
func NewSQLDB(p *Postgres) *sql.DB {
	return stdlib.OpenDBFromPool(p.DB)
}

var Module = fx.Module("database", fx.Provide(NewFx, NewSQLDB))
