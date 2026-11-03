package buyer

import (
	"context"
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

func (s *BuyerService) Register(ctx context.Context, in incoming.RegisterBuyerInput) (buyer.Buyer, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Email) == "" {
		return buyer.Buyer{}, buyer.ErrNotFound
	}
	b := buyer.Buyer{
		Name:    in.Name,
		Email:   in.Email,
		Phone:   in.Phone,
		Address: in.Address,
	}
	return s.repo.Create(ctx, b)
}

func (s *BuyerService) GetByID(ctx context.Context, id int64) (buyer.Buyer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BuyerService) Update(ctx context.Context, id int64, in incoming.UpdateBuyerInput) (buyer.Buyer, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return buyer.Buyer{}, err
	}
	if in.Name != nil {
		b.Name = *in.Name
	}
	if in.Phone != nil {
		b.Phone = *in.Phone
	}
	if in.Address != nil {
		b.Address = *in.Address
	}
	return s.repo.Update(ctx, b)
}

func (s *BuyerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
