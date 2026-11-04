package outgoing

import (
	"Go_Commerce/internal/domain/seller"
	"context"
)

type SellerRepository interface {
	Create(ctx context.Context, s seller.Seller) (seller.Seller, error)
	GetByID(ctx context.Context, id int64) (seller.Seller, error)
	GetByEmail(ctx context.Context, email string) (seller.Seller, error) // GetByEmail returns seller.ErrNotFound when no active seller has this email.
	Update(ctx context.Context, s seller.Seller) (seller.Seller, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]seller.Seller, error)
}
