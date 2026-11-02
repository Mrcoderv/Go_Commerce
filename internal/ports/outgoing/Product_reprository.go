package outgoing

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product duct) error
	GetByID(ctx context.Context, id int64) (Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product Product) error
	Delete(ctx context.Context, id int64) error

}
