package incoming

import (
	"Go_Commerce/internal/domain/seller"
	"context"
)

type RegisterSellerInput struct {
	Name    string
	Email   string
	Phone   string
	Address string
}
type UpdateSellerInput struct {
	Name    *string
	Email   *string
	Phone   *string
	Address *string
}

type SellerService interface {
	// Register creates a seller. Email must be unique among active sellers.
	// Errors: ErrInvalidInput, ErrEmailAlreadyExists (from domain).
	Register(ctx context.Context, in RegisterSellerInput) (seller.Seller, error)
	GetByID(ctx context.Context, id int64) (seller.Seller, error)
	List(ctx context.Context, f ListFilter) ([]seller.Seller, error)
	Update(ctx context.Context, id int64, in UpdateSellerInput) (seller.Seller, error)

	Delete(ctx context.Context, id int64) error
}
