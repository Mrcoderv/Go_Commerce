package incoming

import (
	"Go_Commerce/internal/domain/product"
	"context"
)

type CreateProductInput struct {
	SellerID    int64
	Name        string
	Description string
	Price       int64 // paisa; fixed-point money
	Stock       int
}

type UpdateProductInput struct {
	Name        *string
	Description *string
	Price       *int64 // paisa; fixed-point money
	Stock       *int
}

type ProductFilter struct {
	SellerID *int64
	Limit    int
	Offset   int
}

type ProductService interface {
	Create(ctx context.Context, in CreateProductInput) (product.Product, error)
	GetByID(ctx context.Context, id int64) (product.Product, error)
	List(ctx context.Context, f ProductFilter) ([]product.Product, error)
	Update(ctx context.Context, id int64, sellerID int64, in UpdateProductInput) (product.Product, error)
	Delete(ctx context.Context, id int64, sellerID int64) error
}
