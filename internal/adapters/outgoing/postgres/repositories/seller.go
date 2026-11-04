package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/outgoing"
)

type SellerRepository struct {
	db *sql.DB
}

func NewSellerRepository(db *sql.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

func (r *SellerRepository) Create(ctx context.Context, s seller.Seller) (seller.Seller, error) {
	const q = `INSERT INTO sellers (name, email, phone, address)
	           VALUES ($1, $2, $3, $4) RETURNING id`

	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, s.Name, s.Email, s.Phone, s.Address).
		Scan(&s.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return seller.Seller{}, seller.ErrEmailAlreadyExists
		}
		return seller.Seller{}, fmt.Errorf("create seller: %w", err)
	}
	return s, nil
}

func (r *SellerRepository) GetByID(ctx context.Context, id int64) (seller.Seller, error) {
	const q = `SELECT id, name, email, phone, address
	           FROM sellers WHERE id = $1 AND deleted_at IS NULL`

	var s seller.Seller
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, id).
		Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return seller.Seller{}, seller.ErrNotFound
		}
		return seller.Seller{}, fmt.Errorf("get seller: %w", err)
	}
	return s, nil
}

func (r *SellerRepository) GetByEmail(ctx context.Context, email string) (seller.Seller, error) {
	const q = `SELECT id, name, email, phone, address
	           FROM sellers WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`

	var s seller.Seller
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, email).
		Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return seller.Seller{}, seller.ErrNotFound
		}
		return seller.Seller{}, fmt.Errorf("get seller by email: %w", err)
	}
	return s, nil
}
// Compile-time check: fails the build if a method is missing.
var _ outgoing.SellerRepository = (*SellerRepository)(nil)

func (r *SellerRepository) List(ctx context.Context, limit, offset int) ([]seller.Seller, error) {
	const q = `SELECT id, name, email, phone, address
	           FROM sellers WHERE deleted_at IS NULL
	           ORDER BY id DESC LIMIT $1 OFFSET $2`

	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list sellers: %w", err)
	}
	defer rows.Close()

	sellers := []seller.Seller{}
	for rows.Next() {
		var s seller.Seller
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Address); err != nil {
			return nil, fmt.Errorf("scan seller: %w", err)
		}
		sellers = append(sellers, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list sellers: %w", err)
	}
	return sellers, nil
}

func (r *SellerRepository) Update(ctx context.Context, s seller.Seller) (seller.Seller, error) {
	const q = `UPDATE sellers
	           SET name = $1, email = $2, phone = $3, address = $4, updated_at = NOW()
	           WHERE id = $5 AND deleted_at IS NULL
	           RETURNING id, name, email, phone, address`

	var out seller.Seller
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, s.Name, s.Email, s.Phone, s.Address, s.ID).
		Scan(&out.ID, &out.Name, &out.Email, &out.Phone, &out.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return seller.Seller{}, seller.ErrNotFound
		}
		if isUniqueViolation(err) {
			return seller.Seller{}, seller.ErrEmailAlreadyExists
		}
		return seller.Seller{}, fmt.Errorf("update seller: %w", err)
	}
	return out, nil
}

func (r *SellerRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE sellers SET deleted_at = NOW(), updated_at = NOW()
	           WHERE id = $1 AND deleted_at IS NULL`

	res, err := postgres.Querier(ctx, r.db).ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete seller: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete seller: %w", err)
	}
	if n == 0 {
		return seller.ErrNotFound
	}
	return nil
}