package order

import (
	"context"

	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
)

type OrderService struct {
	repo outgoing.OrderRepository
}

func NewOrderService(repo outgoing.OrderRepository) incoming.OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) PlaceOrder(ctx context.Context, in incoming.PlaceOrderInput) (order.Order, error) {
	return order.Order{}, nil
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) ListByBuyer(ctx context.Context, buyerID int64, f incoming.ListFilter) ([]order.Order, error) {
	return s.repo.ListByBuyer(ctx, buyerID, f.Limit, f.Offset)
}

func (s *OrderService) ListBySeller(ctx context.Context, sellerID int64, f incoming.ListFilter) ([]order.Order, error) {
	return s.repo.ListBySeller(ctx, sellerID, f.Limit, f.Offset)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int64, next order.Status) (order.Order, error) {
	if !next.IsValid() {
		return order.Order{}, order.ErrInvalidInput

	}
	orders, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}
	if !orders.Status.CanTransitionTo(next) {
		return order.Order{}, order.ErrInvalidTransaction
	}
	if err := s.repo.UpdateStatus(ctx, id, next); err != nil {
		return order.Order{}, err
	}
	orders.Status = next
	return orders, nil
}

func (s *OrderService) Cancel(ctx context.Context, id int64) error {
	_, err := s.UpdateStatus(ctx, id, order.StatusCancelled)
	return err
}

var _ incoming.OrderService = (*OrderService)(nil)
