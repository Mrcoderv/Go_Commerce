package seller_test

import (
	"Go_Commerce/internal/core/seller"
	sellerdomain "Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/incoming"
	"context"
	"errors"
	"testing"
)

type fakeSellerRepo struct {
	sellers map[int64]sellerdomain.Seller
	nextID  int64
}

func newFakeSellerRepo() *fakeSellerRepo {
	return &fakeSellerRepo{
		sellers: make(map[int64]sellerdomain.Seller),
		nextID:  1,
	}
}

func (r *fakeSellerRepo) Create(ctx context.Context, s sellerdomain.Seller) (sellerdomain.Seller, error) {
	for _, ex := range r.sellers {
		if ex.Email == s.Email {
			return sellerdomain.Seller{}, sellerdomain.ErrEmailAlreadyExists
		}
	}
	s.ID = r.nextID
	r.nextID++
	r.sellers[s.ID] = s
	return s, nil
}

func (r *fakeSellerRepo) GetByID(ctx context.Context, id int64) (sellerdomain.Seller, error) {
	s, ok := r.sellers[id]
	if !ok {
		return sellerdomain.Seller{}, sellerdomain.ErrNotFound
	}
	return s, nil
}

func (r *fakeSellerRepo) GetByEmail(ctx context.Context, email string) (sellerdomain.Seller, error) {
	for _, s := range r.sellers {
		if s.Email == email {
			return s, nil
		}
	}
	return sellerdomain.Seller{}, sellerdomain.ErrNotFound
}

func (r *fakeSellerRepo) Update(ctx context.Context, s sellerdomain.Seller) (sellerdomain.Seller, error) {
	if _, ok := r.sellers[s.ID]; !ok {
		return sellerdomain.Seller{}, sellerdomain.ErrNotFound
	}
	r.sellers[s.ID] = s
	return s, nil
}

func (r *fakeSellerRepo) Delete(ctx context.Context, id int64) error {
	if _, ok := r.sellers[id]; !ok {
		return sellerdomain.ErrNotFound
	}
	delete(r.sellers, id)
	return nil
}

func (r *fakeSellerRepo) List(ctx context.Context, limit, offset int) ([]sellerdomain.Seller, error) {
	out := make([]sellerdomain.Seller, 0, len(r.sellers))
	for _, s := range r.sellers {
		out = append(out, s)
	}
	return out, nil
}

func TestSellerService_Register(t *testing.T) {
	repo := newFakeSellerRepo()
	svc := seller.NewSellerService(repo)
	ctx := context.Background()

	// 1. Success
	s, err := svc.Register(ctx, incoming.RegisterSellerInput{
		Name:    "Alice",
		Email:   "alice@example.com",
		Phone:   "1234567890",
		Address: "123 Main St",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID != 1 || s.Name != "Alice" {
		t.Errorf("unexpected seller: %+v", s)
	}

	// 2. Duplicate email
	_, err = svc.Register(ctx, incoming.RegisterSellerInput{
		Name:  "Alice Clone",
		Email: "alice@example.com",
	})
	if !errors.Is(err, sellerdomain.ErrEmailAlreadyExists) {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}

	// 3. Validation failure
	_, err = svc.Register(ctx, incoming.RegisterSellerInput{
		Name:  "",
		Email: "bad@example.com",
	})
	if !errors.Is(err, sellerdomain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestSellerService_GetByID_NotFound(t *testing.T) {
	repo := newFakeSellerRepo()
	svc := seller.NewSellerService(repo)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if !errors.Is(err, sellerdomain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
