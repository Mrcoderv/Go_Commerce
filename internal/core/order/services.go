package order

import (
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
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

func (s *OrderService) PlaceOrder(ctx context.Context, in incoming.PlaceOrderInput) (order.Order, error) {
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

			// DecreaseStock exists to prevent the quantity issue
			if err := s.products.DecreaseStock(ctx, it.ProductID, it.Quantity); err != nil {
				return err
			}

			o.Items = append(o.Items, order.Item{
				ProductID:   p.ID,
				SellerID:    p.SellerID,
				ProductName: p.Name,
				UnitPrice:   p.UnitPrice,
				Quantity:    it.Quantity,
			})
			o.Total += p.UnitPrice * int64(it.Quantity)
		}

		//saving order
		created, err := s.repo.Create(ctx, o)
		if err != nil {
			return err
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
	return s.repo.ListByBuyer(ctx, buyerID, f.Limit, f.Offset)
}

func (s *OrderService) ListBySeller(ctx context.Context, sellerID int64, f incoming.ListFilter) ([]order.Order, error) {
	return s.repo.ListBySeller(ctx, sellerID, f.Limit, f.Offset)
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
	_, err := s.UpdateStatus(ctx, id, order.StatusCancelled)
	return err
}

var _ incoming.OrderService = (*OrderService)(nil)

func normalizeItems(in incoming.PlaceOrderInput) ([]incoming.OrderItemInput, error) {
	if len(in.Items) == 0 {
		return nil, order.ErrInvalidInput
	}
	out := make([]incoming.OrderItemInput, 0, len(in.Items))
	for _, it := range in.Items {
		if it.ProductID <= 0 {
			return nil, order.ErrInvalidInput
		}
		if it.Quantity <= 0 {
			return nil, order.ErrInvalidInput
		}
		out = append(out, it)
	}
	return out, nil
}
