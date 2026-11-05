package repositories

import (
	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type BuyerRepository struct {
	db *sql.DB
}

func NewBuyerRepository(db *sql.DB) *BuyerRepository {
	return &BuyerRepository{db: db}
}

var _ outgoing.BuyerRepository = (*BuyerRepository)(nil)

func (r *BuyerRepository) Create(ctx context.Context, b buyer.Buyer) (buyer.Buyer, error) {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}

	query := `INSERT INTO buyers (name, email, phone, address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := postgres.Querier(ctx, r.db).QueryRowContext(ctx, query, b.Name, b.Email, b.Phone, b.Address, b.CreatedAt, now).Scan(&b.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return buyer.Buyer{}, buyer.ErrEmailAlreadyExists
		}
		return buyer.Buyer{}, fmt.Errorf("create buyer: %w", err)
	}

	return b, nil
}

func (r *BuyerRepository) GetByID(ctx context.Context, id int64) (buyer.Buyer, error) {
	b := buyer.Buyer{}
	query := `SELECT id, name, email, phone, address, created_at FROM buyers WHERE id = $1 AND deleted_at IS NULL`
	err := postgres.Querier(ctx, r.db).QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Name, &b.Email, &b.Phone, &b.Address, &b.CreatedAt)
	if err != nil {
		return buyer.Buyer{}, mapBuyerErr(err)
	}
	return b, nil
}

func (r *BuyerRepository) GetByEmail(ctx context.Context, email string) (buyer.Buyer, error) {
	b := buyer.Buyer{}
	query := `SELECT id, name, email, phone, address, created_at FROM buyers WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`
	err := postgres.Querier(ctx, r.db).QueryRowContext(ctx, query, email).Scan(&b.ID, &b.Name, &b.Email, &b.Phone, &b.Address, &b.CreatedAt)
	if err != nil {
		return buyer.Buyer{}, mapBuyerErr(err)
	}
	return b, nil
}

func (r *BuyerRepository) List(ctx context.Context, limit, offset int) ([]buyer.Buyer, error) {
	query := `SELECT id, name, email, phone, address, created_at FROM buyers WHERE deleted_at IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list buyers: %w", err)
	}
	defer rows.Close()

	var buyers []buyer.Buyer
	for rows.Next() {
		var b buyer.Buyer
		if err := rows.Scan(&b.ID, &b.Name, &b.Email, &b.Phone, &b.Address, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan buyer: %w", err)
		}
		buyers = append(buyers, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list buyers: %w", err)
	}
	if buyers == nil {
		buyers = []buyer.Buyer{}
	}
	return buyers, nil
}

func (r *BuyerRepository) Update(ctx context.Context, b buyer.Buyer) (buyer.Buyer, error) {
	query := `UPDATE buyers SET name = $1, email = $2, phone = $3, address = $4, updated_at = NOW() WHERE id = $5 AND deleted_at IS NULL RETURNING id, name, email, phone, address, created_at`
	err := postgres.Querier(ctx, r.db).QueryRowContext(ctx, query, b.Name, b.Email, b.Phone, b.Address, b.ID).Scan(&b.ID, &b.Name, &b.Email, &b.Phone, &b.Address, &b.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return buyer.Buyer{}, buyer.ErrEmailAlreadyExists
		}
		return buyer.Buyer{}, mapBuyerErr(err)
	}
	return b, nil
}

func (r *BuyerRepository) Delete(ctx context.Context, id int64) error {
	res, err := postgres.Querier(ctx, r.db).ExecContext(ctx, "DELETE FROM buyers WHERE id = $1", id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return buyer.ErrConflict
		}
		return fmt.Errorf("delete buyer: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete buyer: %w", err)
	}
	if n == 0 {
		return buyer.ErrNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func mapBuyerErr(er error) error {
	if errors.Is(er, sql.ErrNoRows) {
		return buyer.ErrNotFound
	}
	return er
}
