package seller

import (
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"strings"
)

type SellerService struct {
	repo outgoing.SellerRepository
}

func NewSellerService(repo outgoing.SellerRepository) *SellerService {
	return &SellerService{repo: repo}
}

func (s *SellerService) Register(ctx context.Context, in incoming.RegisterSellerInput) (seller.Seller, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Email) == "" {
		return seller.Seller{}, seller.ErrInvalidInput // or ErrNotFound if you don't have InvalidInput yet
	}
	entity := seller.Seller{
		Name:    in.Name,
		Email:   in.Email,
		Phone:   in.Phone,
		Address: in.Address,
	}
	return s.repo.Create(ctx, entity)
}
