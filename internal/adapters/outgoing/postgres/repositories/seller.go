package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/seller"
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