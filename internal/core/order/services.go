package order

import (
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"fmt"
	"sort"
)

type OrderService struct {
	repo     outgoing.OrderRepository
	buyers   outgoing.BuyerRepository
	products outgoing.ProductRepository
	tx       outgoing.TxManager
}

func NewOrderService(
	repo outgoing.OrderRepository,
	buyers outgoing.BuyerRepository,
	products outgoing.ProductRepository,
	tx outgoing.TxManager,
) incoming.OrderService {
	return &OrderService{repo: repo, buyers: buyers, products: products, tx: tx}
}

var _ incoming.OrderService = (*OrderService)(nil)

func (s *OrderService) PlaceOrder(ctx context.Context, in incoming.PlaceOrderInput) (order.Order, error) {
	if in.BuyerID <= 0 {
		return order.Order{}, order.ErrInvalidInput
	}

	items, err := normalizeItems(in)
	if err != nil {
		return order.Order{}, err
	}

	sort.Slice(items, func(i, j int) bool { return items[i].ProductID < items[j].ProductID })

	var placed order.Order

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// 1. The buyer must exist (buyer.ErrNotFound otherwise).
		if _, err := s.buyers.GetByID(ctx, in.BuyerID); err != nil {
			return err
		}

		o := order.Order{
			BuyerID: in.BuyerID,
			Status:  order.StatusPending,
			Items:   make([]order.Item, 0, len(items)),
		}

		for _, it := range items {
			// 2. Read the product to snapshot its seller, name and price.
			p, err := s.products.GetByID(ctx, it.ProductID)
			if err != nil {
				return err
			}

			// DecreaseStock checks and reduces stock in one statement.
			if err := s.products.DecreaseStock(ctx, it.ProductID, it.Quantity); err != nil {
				return err
			}

			o.Items = append(o.Items, order.Item{
				ProductID:   p.ID,
				SellerID:    p.SellerID,
				ProductName: p.Name,
				UnitPrice:   int64(p.Price),
				Quantity:    it.Quantity,
			})
			o.Total += int64(p.Price * float64(it.Quantity))
		}

		// Save order
		created, err := s.repo.Create(ctx, o)
		if err != nil {
			return fmt.Errorf("save order: %w", err)
		}
		placed = created
		return nil
	})
	if err != nil {
		return order.Order{}, err
	}
	return placed, nil
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) ListByBuyer(ctx context.Context, buyerID int64, f incoming.ListFilter) ([]order.Order, error) {
	limit, offset := clampFilter(f)
	return s.repo.ListByBuyer(ctx, buyerID, limit, offset)
}

func (s *OrderService) ListBySeller(ctx context.Context, sellerID int64, f incoming.ListFilter) ([]order.Order, error) {
	limit, offset := clampFilter(f)
	return s.repo.ListBySeller(ctx, sellerID, limit, offset)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int64, next order.Status) (order.Order, error) {
	if !next.IsValid() {
		return order.Order{}, order.ErrInvalidInput
	}

	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}
	if !o.Status.CanTransitionTo(next) {
		return order.Order{}, order.ErrInvalidTransition
	}
	if err := s.repo.UpdateStatus(ctx, id, next); err != nil {
		return order.Order{}, err
	}

	o.Status = next
	return o, nil
}

func (s *OrderService) Cancel(ctx context.Context, id int64) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		o, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if !o.Status.CanTransitionTo(order.StatusCancelled) {
			return order.ErrInvalidTransition
		}
		for _, it := range o.Items {
			if err := s.products.IncreaseStock(ctx, it.ProductID, it.Quantity); err != nil {
				return err
			}
		}
		return s.repo.UpdateStatus(ctx, id, order.StatusCancelled)
	})
}

func clampFilter(f incoming.ListFilter) (int, int) {
	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func normalizeItems(in incoming.PlaceOrderInput) ([]incoming.OrderItemInput, error) {
	if len(in.Items) == 0 {
		return nil, order.ErrInvalidInput
	}
	merged := make(map[int64]int)
	var orderKey []int64
	for _, it := range in.Items {
		if it.ProductID <= 0 || it.Quantity <= 0 {
			return nil, order.ErrInvalidInput
		}
		if _, seen := merged[it.ProductID]; !seen {
			orderKey = append(orderKey, it.ProductID)
		}
		merged[it.ProductID] += it.Quantity
	}

	out := make([]incoming.OrderItemInput, 0, len(merged))
	for _, pid := range orderKey {
		out = append(out, incoming.OrderItemInput{
			ProductID: pid,
			Quantity:  merged[pid],
		})
	}
	return out, nil
}
