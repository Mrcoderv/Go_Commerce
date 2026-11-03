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

func NewOrderService(repo outgoing.OrderRepository) *OrderService {
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
	if err := s.repo.UpdateStatus(ctx, id, next); err != nil {
		return order.Order{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) Cancel(ctx context.Context, id int64) error {
	return s.repo.UpdateStatus(ctx, id, order.Status(""))
}
