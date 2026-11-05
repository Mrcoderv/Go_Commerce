package outgoing

import (
	"Go_Commerce/internal/domain/buyer"
	"context"
)

type BuyerRepository interface {
	Create(ctx context.Context, b buyer.Buyer) (buyer.Buyer, error)
	GetByID(ctx context.Context, id int64) (buyer.Buyer, error)
	GetByEmail(ctx context.Context, email string) (buyer.Buyer, error)
	List(ctx context.Context, limit, offset int) ([]buyer.Buyer, error)
	Update(ctx context.Context, b buyer.Buyer) (buyer.Buyer, error)
	Delete(ctx context.Context, id int64) error
}
