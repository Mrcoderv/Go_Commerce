package buyer

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
)

type BuyerService struct {
	repo outgoing.BuyerRepository
	log  logger.Logger
}

func NewBuyerService(repo outgoing.BuyerRepository, log logger.Logger) incoming.BuyerService {
	return &BuyerService{repo: repo, log: log}
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

	created, err := s.repo.Create(ctx, newBuyer)
	if err != nil {
		s.log.Error("buyer registration failed", "email", email, "error", err)
		return buyer.Buyer{}, err
	}
	s.log.Info("buyer registered", "buyer_id", created.ID)
	return created, nil
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
	updated, err := s.repo.Update(ctx, b)
	if err != nil {
		s.log.Error("buyer update failed", "buyer_id", id, "error", err)
		return buyer.Buyer{}, err
	}
	s.log.Info("buyer updated", "buyer_id", id)
	return updated, nil
}

func (s *BuyerService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Info("buyer deleted", "buyer_id", id)
	return nil
}
