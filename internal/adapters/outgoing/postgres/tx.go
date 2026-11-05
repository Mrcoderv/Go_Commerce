package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Go_Commerce/internal/ports/outgoing"
)

type txKey struct{}

// DBTX is what every repository runs its queries on.
// Both *sql.DB and *sql.Tx have these three methods.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Querier returns the transaction stored in ctx, or the plain db if there is none.
func Querier(ctx context.Context, db *sql.DB) DBTX {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) outgoing.TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Already inside a transaction: join it instead of starting another.
	if _, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			return fmt.Errorf("%w (rollback failed: %v)", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
