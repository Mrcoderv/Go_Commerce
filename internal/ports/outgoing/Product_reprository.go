package outgoing

import (
	"Go_Commerce/internal/domain/product"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product product.Product) error
	GetByID(ctx context.Context, id int64) (product.Product, error)
	GetAll(ctx context.Context) ([]product.Product, error)
	Update(ctx context.Context, product product.Product) error
	Delete(ctx context.Context, id int64) error
}
