package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/ports/outgoing"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ outgoing.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u auth.User) (auth.User, error) {
	const q = `INSERT INTO users (email, password_hash, role, seller_id, buyer_id)
	           VALUES ($1, $2, $3, $4, $5)
	           RETURNING id, created_at`

	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, u.Email, u.PasswordHash, u.Role, u.SellerID, u.BuyerID).
		Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return auth.User{}, auth.ErrEmailExists
		}
		return auth.User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByEmailAndRole(ctx context.Context, email string, role auth.Role) (auth.User, error) {
	const q = `SELECT id, email, password_hash, role, seller_id, buyer_id, created_at
	           FROM users WHERE LOWER(email)=LOWER($1) AND role=$2`

	var u auth.User
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, email, role).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.SellerID, &u.BuyerID, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (auth.User, error) {
	const q = `SELECT id, email, password_hash, role, seller_id, buyer_id, created_at
	           FROM users WHERE id=$1`
	var u auth.User
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.SellerID, &u.BuyerID, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}
