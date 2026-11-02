package incoming

import (
	"Go_Commerce/internal/domain/buyer"
	"context"
)

type BuyerService interface {
	Register(ctx context.Context, in RegisterBuyerInput) (buyer.Buyer, error)
	GetByID(ctx context.Context, id int64) (buyer.Buyer, error)
	Update(ctx context.Context, id int64, in UpdateBuyerInput) (buyer.Buyer, error)
	Delete(ctx context.Context, id int64) error
}
