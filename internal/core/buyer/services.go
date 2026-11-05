package buyer

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
)

type BuyerService struct {
	repo outgoing.BuyerRepository
}

func NewBuyerService(repo outgoing.BuyerRepository) incoming.BuyerService {
	return &BuyerService{repo: repo}
}

var _ incoming.BuyerService = (*BuyerService)(nil)

func (s *BuyerService) Register(ctx context.Context, in incoming.RegisterBuyerInput) (buyer.Buyer, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))

	if name == "" || email == "" {
		return buyer.Buyer{}, buyer.ErrInvalidInput
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return buyer.Buyer{}, buyer.ErrInvalidInput
	}

	existingBuyer, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingBuyer.ID != 0 {
		return buyer.Buyer{}, buyer.ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, buyer.ErrNotFound) {
		return buyer.Buyer{}, fmt.Errorf("check email: %w", err)
	}

	newBuyer := buyer.Buyer{
		Name:    name,
		Email:   email,
		Phone:   strings.TrimSpace(in.Phone),
		Address: strings.TrimSpace(in.Address),
	}

	return s.repo.Create(ctx, newBuyer)
}

func (s *BuyerService) GetByID(ctx context.Context, id int64) (buyer.Buyer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BuyerService) List(ctx context.Context, f incoming.ListFilter) ([]buyer.Buyer, error) {
	limit, offset := f.Limit, f.Offset
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *BuyerService) Update(ctx context.Context, id int64, in incoming.UpdateBuyerInput) (buyer.Buyer, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return buyer.Buyer{}, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return buyer.Buyer{}, buyer.ErrInvalidInput
		}
		b.Name = name
	}
	if in.Phone != nil {
		b.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Address != nil {
		b.Address = strings.TrimSpace(*in.Address)
	}
	return s.repo.Update(ctx, b)
}

func (s *BuyerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
