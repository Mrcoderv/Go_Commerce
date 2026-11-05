package product_test

import (
	"Go_Commerce/internal/core/product"
	productdomain "Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"
	"context"
	"errors"
	"testing"
)

type fakeProductRepo struct {
	products map[int64]productdomain.Product
	nextID   int64
}

func newFakeProductRepo() *fakeProductRepo {
	return &fakeProductRepo{
		products: make(map[int64]productdomain.Product),
		nextID:   1,
	}
}

func (r *fakeProductRepo) Create(ctx context.Context, p productdomain.Product) (productdomain.Product, error) {
	p.ID = r.nextID
	r.nextID++
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
	out := make([]productdomain.Product, 0)
	for _, p := range r.products {
		if sellerID == nil || p.SellerID == *sellerID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *fakeProductRepo) Update(ctx context.Context, p productdomain.Product) (productdomain.Product, error) {
	if _, ok := r.products[p.ID]; !ok {
		return productdomain.Product{}, productdomain.ErrNotFound
	}
	r.products[p.ID] = p
	return p, nil
}

func (r *fakeProductRepo) Delete(ctx context.Context, id int64) error {
	if _, ok := r.products[id]; !ok {
		return productdomain.ErrNotFound
	}
	delete(r.products, id)
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

func TestProductService_Create(t *testing.T) {
	repo := newFakeProductRepo()
	svc := product.NewProductService(repo)
	ctx := context.Background()

	// 1. Success
	p, err := svc.Create(ctx, incoming.CreateProductInput{
		SellerID:    1,
		Name:        "Laptop",
		Description: "A powerful laptop",
		Price:       99900,
		Stock:       10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 || p.Price != 99900 {
		t.Errorf("unexpected product: %+v", p)
	}

	// 2. Validation failures
	tests := []struct {
		name string
		in   incoming.CreateProductInput
	}{
		{"empty name", incoming.CreateProductInput{SellerID: 1, Name: "", Price: 100, Stock: 1}},
		{"zero sellerID", incoming.CreateProductInput{SellerID: 0, Name: "Book", Price: 100, Stock: 1}},
		{"zero price", incoming.CreateProductInput{SellerID: 1, Name: "Book", Price: 0, Stock: 1}},
		{"negative stock", incoming.CreateProductInput{SellerID: 1, Name: "Book", Price: 100, Stock: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(ctx, tt.in)
			if !errors.Is(err, productdomain.ErrInvalidInput) {
				t.Errorf("expected ErrInvalidInput, got %v", err)
			}
		})
	}
}

func TestProductService_OwnershipForbidden(t *testing.T) {
	repo := newFakeProductRepo()
	svc := product.NewProductService(repo)
	ctx := context.Background()

	p, err := svc.Create(ctx, incoming.CreateProductInput{
		SellerID: 1,
		Name:     "Keyboard",
		Price:    5000,
		Stock:    5,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Another seller tries to update
	newName := "Hacked Keyboard"
	_, err = svc.Update(ctx, p.ID, 2, incoming.UpdateProductInput{Name: &newName})
	if !errors.Is(err, productdomain.ErrForbidden) {
		t.Errorf("expected ErrForbidden on update from non-owner, got %v", err)
	}

	// Another seller tries to delete
	err = svc.Delete(ctx, p.ID, 2)
	if !errors.Is(err, productdomain.ErrForbidden) {
		t.Errorf("expected ErrForbidden on delete from non-owner, got %v", err)
	}

	// Owner can update
	updated, err := svc.Update(ctx, p.ID, 1, incoming.UpdateProductInput{Name: &newName})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Hacked Keyboard" {
		t.Errorf("expected name to be updated, got %s", updated.Name)
	}
}

func TestProductService_StockOperations(t *testing.T) {
	repo := newFakeProductRepo()
	svc := product.NewProductService(repo)
	ctx := context.Background()

	p, _ := svc.Create(ctx, incoming.CreateProductInput{
		SellerID: 1,
		Name:     "Mouse",
		Price:    2000,
		Stock:    2,
	})

	// Check availability
	_, err := svc.CheckAvailability(ctx, p.ID, 3)
	if !errors.Is(err, productdomain.ErrInsufficientStock) {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}

	// Reduce stock successfully
	err = svc.ReduceStock(ctx, p.ID, 2)
	if err != nil {
		t.Fatalf("reduce stock failed: %v", err)
	}

	// Reduce stock now fails
	err = svc.ReduceStock(ctx, p.ID, 1)
	if !errors.Is(err, productdomain.ErrInsufficientStock) {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}
