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
	ProductID int64
	BuyerID   int64
	Items     []OrderItemInput
}

type OrderService interface {
	// PlaceOrder is the main flow, all inside one transaction:
	// validate buyer -> load products -> check stock -> snapshot prices
	// -> decrease stock -> save order and items.
	// Errors: ErrInvalidInput (empty cart, qty <= 0), buyer.ErrNotFound,
	// product.ErrNotFound, product.ErrInsufficientStock.
	PlaceOrder(ctx context.Context, in PlaceOrderInput) (order.Order, error)

	// GetByID returns the order with its items. Errors: ErrNotFound.
	GetByID(ctx context.Context, id int64) (order.Order, error)

	// ListByBuyer is the buyer's history, newest first.
	ListByBuyer(ctx context.Context, buyerID int64, f ListFilter) ([]order.Order, error)

	// ListBySeller returns orders containing this seller's items.
	ListBySeller(ctx context.Context, sellerID int64, f ListFilter) ([]order.Order, error)

	// UpdateStatus enforces legal transitions via Status.CanTransitionTo.
	// Errors: ErrNotFound, ErrInvalidTransition.
	UpdateStatus(ctx context.Context, id int64, next order.Status) (order.Order, error)

	// Cancel sets the status to cancelled and restores stock, in one transaction.
	// Errors: ErrNotFound, ErrInvalidTransition (already shipped).
	Cancel(ctx context.Context, id int64) error
}
