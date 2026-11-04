package outgoing

import (
	"Go_Commerce/internal/domain/product"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product product.Product) (product.Product, error)
	GetByID(ctx context.Context, id int64) (product.Product, error)
	GetAll(ctx context.Context) ([]product.Product, error)
	Update(ctx context.Context, product product.Product) (product.Product, error)
	Delete(ctx context.Context, id int64) error
	DecreasesStock(ctx context.Context, id product.Product, qty product.Product) error
	//	IncreasesStock(ctx context.Context, id product.Product, qty product.Product) error
}
