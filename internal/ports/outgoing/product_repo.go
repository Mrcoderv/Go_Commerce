package outgoing

import (
	"Go_Commerce/internal/domain/product"
	"context"
)

type ProductFilter struct {
	SellerID *int64
	Limit    int
	Offset   int
}

type ProductRepository interface {
	Create(ctx context.Context, p product.Product) (product.Product, error)
	GetByID(ctx context.Context, id int64) (product.Product, error)
	List(ctx context.Context, sellerID *int64, limit, offset int) ([]product.Product, error)
	Update(ctx context.Context, p product.Product) (product.Product, error)
	Delete(ctx context.Context, id int64) error
	DecreaseStock(ctx context.Context, id int64, qty int) error
	IncreaseStock(ctx context.Context, id int64, qty int) error
}
