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

func (r *SellerRepo) GetByEmail(ctx context.Context, email string) (seller.Seller, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.byID {
		if existing.Email == email {
			return existing, nil
		}
	}
	return seller.Seller{}, seller.ErrNotFound
}

func (r *SellerRepo) Update(ctx context.Context, s seller.Seller) (seller.Seller, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.byID[s.ID]
	if !ok {
		return seller.Seller{}, seller.ErrNotFound
	}
	if s.Name != "" {
		existing.Name = s.Name
	}
	if s.Email != "" {
		for id, ex := range r.byID {
			if id != s.ID && ex.Email == s.Email {
				return seller.Seller{}, seller.ErrEmailAlreadyExists
			}
		}
		existing.Email = s.Email
	}
	if s.Phone != "" {
		existing.Phone = s.Phone
	}
	if s.Address != "" {
		existing.Address = s.Address
	}
	r.byID[s.ID] = existing
	return existing, nil
}

func (r *SellerRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return seller.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

func (r *SellerRepo) List(ctx context.Context, limit, offset int) ([]seller.Seller, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]seller.Seller, 0)
	for _, s := range r.byID {
		res = append(res, s)
	}
	if offset >= len(res) {
		return []seller.Seller{}, nil
	}
	end := offset + limit
	if end > len(res) {
		end = len(res)
	}
	return res[offset:end], nil
}
