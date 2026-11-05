package order_test

import (
	"Go_Commerce/internal/core/order"
	buyerdomain "Go_Commerce/internal/domain/buyer"
	orderdomain "Go_Commerce/internal/domain/order"
	productdomain "Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"
	"context"
	"errors"
	"testing"
)

// --- In-memory Fakes ---

type fakeBuyerRepo struct {
	buyers map[int64]buyerdomain.Buyer
}

func (r *fakeBuyerRepo) Create(ctx context.Context, b buyerdomain.Buyer) (buyerdomain.Buyer, error) {
	r.buyers[b.ID] = b
	return b, nil
}
func (r *fakeBuyerRepo) GetByID(ctx context.Context, id int64) (buyerdomain.Buyer, error) {
	b, ok := r.buyers[id]
	if !ok {
		return buyerdomain.Buyer{}, buyerdomain.ErrNotFound
	}
	return b, nil
}
func (r *fakeBuyerRepo) GetByEmail(ctx context.Context, email string) (buyerdomain.Buyer, error) {
	return buyerdomain.Buyer{}, nil
}
func (r *fakeBuyerRepo) List(ctx context.Context, limit, offset int) ([]buyerdomain.Buyer, error) {
	return nil, nil
}
func (r *fakeBuyerRepo) Update(ctx context.Context, b buyerdomain.Buyer) (buyerdomain.Buyer, error) {
	return b, nil
}
func (r *fakeBuyerRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

type fakeProductRepo struct {
	products map[int64]productdomain.Product
}

func (r *fakeProductRepo) Create(ctx context.Context, p productdomain.Product) (productdomain.Product, error) {
	r.products[p.ID] = p
	return p, nil
}
func (r *fakeProductRepo) GetByID(ctx context.Context, id int64) (productdomain.Product, error) {
	p, ok := r.products[id]
	if !ok {
		return productdomain.Product{}, productdomain.ErrNotFound
	}
	return p, nil
}
func (r *fakeProductRepo) List(ctx context.Context, sellerID *int64, limit, offset int) ([]productdomain.Product, error) {
	return nil, nil
}
func (r *fakeProductRepo) Update(ctx context.Context, p productdomain.Product) (productdomain.Product, error) {
	return p, nil
}
func (r *fakeProductRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
func (r *fakeProductRepo) DecreaseStock(ctx context.Context, id int64, qty int) error {
	p, ok := r.products[id]
	if !ok {
		return productdomain.ErrNotFound
	}
	if p.Stock < qty {
		return productdomain.ErrInsufficientStock
	}
	p.Stock -= qty
	r.products[id] = p
	return nil
}
func (r *fakeProductRepo) IncreaseStock(ctx context.Context, id int64, qty int) error {
	p, ok := r.products[id]
	if !ok {
		return productdomain.ErrNotFound
	}
	p.Stock += qty
	r.products[id] = p
	return nil
}

type fakeOrderRepo struct {
	orders map[int64]orderdomain.Order
	nextID int64
}

func (r *fakeOrderRepo) Create(ctx context.Context, o orderdomain.Order) (orderdomain.Order, error) {
	o.ID = r.nextID
	r.nextID++
	r.orders[o.ID] = o
	return o, nil
}
func (r *fakeOrderRepo) GetByID(ctx context.Context, id int64) (orderdomain.Order, error) {
	o, ok := r.orders[id]
	if !ok {
		return orderdomain.Order{}, orderdomain.ErrNotFound
	}
	return o, nil
}
func (r *fakeOrderRepo) ListByBuyer(ctx context.Context, buyerID int64, limit, offset int) ([]orderdomain.Order, error) {
	var out []orderdomain.Order
	for _, o := range r.orders {
		if o.BuyerID == buyerID {
			out = append(out, o)
		}
	}
	return out, nil
}
func (r *fakeOrderRepo) ListBySeller(ctx context.Context, sellerID int64, limit, offset int) ([]orderdomain.Order, error) {
	var out []orderdomain.Order
	for _, o := range r.orders {
		for _, item := range o.Items {
			if item.SellerID == sellerID {
				out = append(out, o)
				break
			}
		}
	}
	return out, nil
}
func (r *fakeOrderRepo) UpdateStatus(ctx context.Context, id int64, status orderdomain.Status) error {
	o, ok := r.orders[id]
	if !ok {
		return orderdomain.ErrNotFound
	}
	o.Status = status
	r.orders[id] = o
	return nil
}

type fakeTxManager struct {
	productRepo *fakeProductRepo
}

func (m *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Snapshot product stock
	snapshot := make(map[int64]int)
	for id, p := range m.productRepo.products {
		snapshot[id] = p.Stock
	}

	err := fn(ctx)
	if err != nil {
		// Rollback product stock to snapshot
		for id, stock := range snapshot {
			p := m.productRepo.products[id]
			p.Stock = stock
			m.productRepo.products[id] = p
		}
		return err
	}
	return nil
}

// --- Tests ---

func setupOrderService() (*fakeBuyerRepo, *fakeProductRepo, *fakeOrderRepo, incoming.OrderService) {
	buyers := &fakeBuyerRepo{buyers: make(map[int64]buyerdomain.Buyer)}
	products := &fakeProductRepo{products: make(map[int64]productdomain.Product)}
	orders := &fakeOrderRepo{orders: make(map[int64]orderdomain.Order), nextID: 1}
	tx := &fakeTxManager{productRepo: products}

	svc := order.NewOrderService(orders, buyers, products, tx)
	return buyers, products, orders, svc
}

func TestOrderService_PlaceOrder_Success(t *testing.T) {
	buyers, products, _, svc := setupOrderService()
	ctx := context.Background()

	buyers.buyers[1] = buyerdomain.Buyer{ID: 1, Name: "Buyer 1"}
	products.products[10] = productdomain.Product{ID: 10, SellerID: 5, Name: "Item 1", Price: 1000, Stock: 5}

	o, err := svc.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: 1,
		Items:   []incoming.OrderItemInput{{ProductID: 10, Quantity: 2}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o.Total != 2000 {
		t.Errorf("expected total 2000, got %d", o.Total)
	}
	if o.Status != orderdomain.StatusPending {
		t.Errorf("expected status pending, got %s", o.Status)
	}
	if products.products[10].Stock != 3 {
		t.Errorf("expected product stock to decrease to 3, got %d", products.products[10].Stock)
	}
}

func TestOrderService_PlaceOrder_ValidationFailures(t *testing.T) {
	buyers, products, _, svc := setupOrderService()
	ctx := context.Background()

	buyers.buyers[1] = buyerdomain.Buyer{ID: 1, Name: "Buyer 1"}
	products.products[10] = productdomain.Product{ID: 10, Price: 1000, Stock: 5}

	// 1. Empty items
	_, err := svc.PlaceOrder(ctx, incoming.PlaceOrderInput{BuyerID: 1, Items: nil})
	if !errors.Is(err, orderdomain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for empty items, got %v", err)
	}

	// 2. Buyer ID <= 0
	_, err = svc.PlaceOrder(ctx, incoming.PlaceOrderInput{BuyerID: 0, Items: []incoming.OrderItemInput{{ProductID: 10, Quantity: 1}}})
	if !errors.Is(err, orderdomain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for non-positive BuyerID, got %v", err)
	}

	// 3. Buyer not found
	_, err = svc.PlaceOrder(ctx, incoming.PlaceOrderInput{BuyerID: 999, Items: []incoming.OrderItemInput{{ProductID: 10, Quantity: 1}}})
	if !errors.Is(err, buyerdomain.ErrNotFound) {
		t.Errorf("expected buyer.ErrNotFound, got %v", err)
	}
}

func TestOrderService_PlaceOrder_RollbackOnSecondItemFailure(t *testing.T) {
	buyers, products, _, svc := setupOrderService()
	ctx := context.Background()

	buyers.buyers[1] = buyerdomain.Buyer{ID: 1, Name: "Buyer 1"}
	products.products[1] = productdomain.Product{ID: 1, SellerID: 5, Name: "Item 1", Price: 1000, Stock: 10}
	products.products[2] = productdomain.Product{ID: 2, SellerID: 5, Name: "Item 2", Price: 2000, Stock: 1} // only 1 in stock

	// Attempt to order 2 units of Item 1 (has 10) and 5 units of Item 2 (has only 1)
	_, err := svc.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: 1,
		Items: []incoming.OrderItemInput{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 5},
		},
	})
	if !errors.Is(err, productdomain.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}

	// CRITICAL CHECK: Product 1's stock must NOT be reduced (rolled back to 10)!
	if products.products[1].Stock != 10 {
		t.Errorf("expected Product 1 stock to roll back to 10, got %d", products.products[1].Stock)
	}
	if products.products[2].Stock != 1 {
		t.Errorf("expected Product 2 stock to remain 1, got %d", products.products[2].Stock)
	}
}

func TestOrderService_UpdateStatus_Transitions(t *testing.T) {
	buyers, products, _, svc := setupOrderService()
	ctx := context.Background()

	buyers.buyers[1] = buyerdomain.Buyer{ID: 1, Name: "Buyer 1"}
	products.products[1] = productdomain.Product{ID: 1, SellerID: 5, Name: "Item 1", Price: 1000, Stock: 10}

	o, _ := svc.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: 1,
		Items:   []incoming.OrderItemInput{{ProductID: 1, Quantity: 1}},
	})

	// Legal: pending -> confirmed
	updated, err := svc.UpdateStatus(ctx, o.ID, orderdomain.StatusConfirmed)
	if err != nil {
		t.Fatalf("legal transition failed: %v", err)
	}
	if updated.Status != orderdomain.StatusConfirmed {
		t.Errorf("expected status confirmed, got %s", updated.Status)
	}

	// Illegal: confirmed -> delivered (must go through shipped)
	_, err = svc.UpdateStatus(ctx, o.ID, orderdomain.StatusDelivered)
	if !errors.Is(err, orderdomain.ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}

	// Legal: confirmed -> shipped
	updated, err = svc.UpdateStatus(ctx, o.ID, orderdomain.StatusShipped)
	if err != nil {
		t.Fatalf("legal transition failed: %v", err)
	}
	if updated.Status != orderdomain.StatusShipped {
		t.Errorf("expected status shipped, got %s", updated.Status)
	}

	// Legal: shipped -> delivered
	updated, err = svc.UpdateStatus(ctx, o.ID, orderdomain.StatusDelivered)
	if err != nil {
		t.Fatalf("legal transition failed: %v", err)
	}
	if updated.Status != orderdomain.StatusDelivered {
		t.Errorf("expected status delivered, got %s", updated.Status)
	}
}

func TestOrderService_CancelRestoresStock(t *testing.T) {
	buyers, products, _, svc := setupOrderService()
	ctx := context.Background()

	buyers.buyers[1] = buyerdomain.Buyer{ID: 1, Name: "Buyer 1"}
	products.products[1] = productdomain.Product{ID: 1, SellerID: 5, Name: "Item 1", Price: 1000, Stock: 5}

	o, _ := svc.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: 1,
		Items:   []incoming.OrderItemInput{{ProductID: 1, Quantity: 2}},
	})
	if products.products[1].Stock != 3 {
		t.Fatalf("expected stock 3 after place order, got %d", products.products[1].Stock)
	}

	// Cancel order
	err := svc.Cancel(ctx, o.ID)
	if err != nil {
		t.Fatalf("cancel order failed: %v", err)
	}

	// Stock restored to 5!
	if products.products[1].Stock != 5 {
		t.Errorf("expected stock to be restored to 5 after cancel, got %d", products.products[1].Stock)
	}

	// Cancelling again is an illegal transition (cancelled -> cancelled)
	err = svc.Cancel(ctx, o.ID)
	if !errors.Is(err, orderdomain.ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition on re-cancelling, got %v", err)
	}
}
