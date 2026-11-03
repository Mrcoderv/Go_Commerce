package outgoing

import (
	"Go_Commerce/internal/domain/order"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, o order.Order) (order.Order, error)
	GetByID(ctx context.Context, id int64) (order.Order, error)
	ListByBuyer(ctx context.Context, buyerID int64, limit, offset int) ([]order.Order, error)
	ListBySeller(ctx context.Context, sellerID int64, limit, offset int) ([]order.Order, error)
	UpdateStatus(ctx context.Context, id int64, status order.Status) error
}
