// ports/incoming/buyer.go
package incoming

import (
	"Go_Commerce/internal/domain/buyer"
	"context"
)

type RegisterBuyerInput struct {
	Name, Email, Phone, Address string
}

// Pointer fields mean "only change what was sent".
type UpdateBuyerInput struct {
	Name, Phone, Address *string
}

type BuyerService interface {
	// Register validates input and creates a buyer.
	// Errors: ErrInvalidInput, ErrEmailAlreadyExists.
	Register(ctx context.Context, in RegisterBuyerInput) (buyer.Buyer, error)

	// GetByID errors: ErrNotFound.
	GetByID(ctx context.Context, id int64) (buyer.Buyer, error)

	// List returns paginated buyers.
	List(ctx context.Context, f ListFilter) ([]buyer.Buyer, error)

	// Update errors: ErrNotFound, ErrInvalidInput.
	Update(ctx context.Context, id int64, in UpdateBuyerInput) (buyer.Buyer, error)

	// Delete errors: ErrNotFound, ErrConflict.
	Delete(ctx context.Context, id int64) error
}
