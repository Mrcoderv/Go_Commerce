package buyer_test

import (
	"Go_Commerce/internal/core/buyer"
	buyerdomain "Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"context"
	"errors"
	"testing"
)

type fakeBuyerRepo struct {
	buyers map[int64]buyerdomain.Buyer
	nextID int64
}

func newFakeBuyerRepo() *fakeBuyerRepo {
	return &fakeBuyerRepo{
		buyers: make(map[int64]buyerdomain.Buyer),
		nextID: 1,
	}
}

func (r *fakeBuyerRepo) Create(ctx context.Context, b buyerdomain.Buyer) (buyerdomain.Buyer, error) {
	for _, ex := range r.buyers {
		if ex.Email == b.Email {
			return buyerdomain.Buyer{}, buyerdomain.ErrEmailAlreadyExists
		}
	}
	b.ID = r.nextID
	r.nextID++
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
	for _, b := range r.buyers {
		if b.Email == email {
			return b, nil
		}
	}
	return buyerdomain.Buyer{}, buyerdomain.ErrNotFound
}

func (r *fakeBuyerRepo) List(ctx context.Context, limit, offset int) ([]buyerdomain.Buyer, error) {
	out := make([]buyerdomain.Buyer, 0, len(r.buyers))
	for _, b := range r.buyers {
		out = append(out, b)
	}
	return out, nil
}

func (r *fakeBuyerRepo) Update(ctx context.Context, b buyerdomain.Buyer) (buyerdomain.Buyer, error) {
	if _, ok := r.buyers[b.ID]; !ok {
		return buyerdomain.Buyer{}, buyerdomain.ErrNotFound
	}
	r.buyers[b.ID] = b
	return b, nil
}

func (r *fakeBuyerRepo) Delete(ctx context.Context, id int64) error {
	if _, ok := r.buyers[id]; !ok {
		return buyerdomain.ErrNotFound
	}
	delete(r.buyers, id)
	return nil
}

func TestBuyerService_Register(t *testing.T) {
	repo := newFakeBuyerRepo()
	svc := buyer.NewBuyerService(repo, logger.NewNop())
	ctx := context.Background()

	// 1. Success
	b, err := svc.Register(ctx, incoming.RegisterBuyerInput{
		Name:    "Bob",
		Email:   "bob@example.com",
		Phone:   "9876543210",
		Address: "456 Market St",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != 1 || b.Name != "Bob" {
		t.Errorf("unexpected buyer: %+v", b)
	}

	// 2. Duplicate email
	_, err = svc.Register(ctx, incoming.RegisterBuyerInput{
		Name:  "Bob Duplicate",
		Email: "bob@example.com",
	})
	if !errors.Is(err, buyerdomain.ErrEmailAlreadyExists) {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}

	// 3. Invalid email / input
	_, err = svc.Register(ctx, incoming.RegisterBuyerInput{
		Name:  "Bob",
		Email: "not-an-email",
	})
	if !errors.Is(err, buyerdomain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestBuyerService_GetByID_NotFound(t *testing.T) {
	repo := newFakeBuyerRepo()
	svc := buyer.NewBuyerService(repo, logger.NewNop())
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 404)
	if !errors.Is(err, buyerdomain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
