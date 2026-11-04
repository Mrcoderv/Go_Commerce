package incoming

import (
	"Go_Commerce/internal/domain/order"
	"context"
)

type OrderItemInput struct {
	ProductID int64
	Quantity  int
}

type PlaceOrderInput struct {
	BuyerID int64
	Items   []OrderItemInput
}

//type ProductReader struct {
//	ID       int64
//	SellerID int64
//	Name     string
//	Price    float64
//	Stock    int
//}

type OrderService interface {
	PlaceOrder(ctx context.Context, in PlaceOrderInput) (order.Order, error)

	GetByID(ctx context.Context, id int64) (order.Order, error)

	ListByBuyer(ctx context.Context, buyerID int64, f ListFilter) ([]order.Order, error)

	ListBySeller(ctx context.Context, sellerID int64, f ListFilter) ([]order.Order, error)
	UpdateStatus(ctx context.Context, id int64, next order.Status) (order.Order, error)
	Cancel(ctx context.Context, id int64) error
}
