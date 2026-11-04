package memory

import (
	"context"
	"sync"

	"Go_Commerce/internal/domain/seller"
)

type SellerRepo struct {
	mu     sync.Mutex
	nextID int64
	byID   map[int64]seller.Seller
}

func NewSellerRepo() *SellerRepo {
	return &SellerRepo{
		nextID: 1,
		byID:   make(map[int64]seller.Seller),
	}
}

func (r *SellerRepo) Create(ctx context.Context, s seller.Seller) (seller.Seller, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.byID {
		if existing.Email == s.Email {
			return seller.Seller{}, seller.ErrEmailAlreadyExists
		}
	}

	s.ID = r.nextID
	r.nextID++
	r.byID[s.ID] = s
	return s, nil
}

func (r *SellerRepo) GetByID(ctx context.Context, id int64) (seller.Seller, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.byID[id]
	if !ok {
		return seller.Seller{}, seller.ErrNotFound
	}
	return s, nil
}

// GetByEmail, Update, Delete, List — same idea (mutex + map)
