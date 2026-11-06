package seller

import (
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"errors"
	"strings"
)

type SellerService struct {
	repo outgoing.SellerRepository
	log  logger.Logger
}

func NewSellerService(repo outgoing.SellerRepository, log logger.Logger) incoming.SellerService {
	return &SellerService{repo: repo, log: log}
}

func (s *SellerService) Register(ctx context.Context, in incoming.RegisterSellerInput) (seller.Seller, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if name == "" || email == "" {
		return seller.Seller{}, seller.ErrInvalidInput
	}

	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return seller.Seller{}, seller.ErrEmailAlreadyExists
	}
	if !errors.Is(err, seller.ErrNotFound) {
		return seller.Seller{}, err
	}

	created, err := s.repo.Create(ctx, seller.Seller{
		Name:    name,
		Email:   email,
		Phone:   strings.TrimSpace(in.Phone),
		Address: strings.TrimSpace(in.Address),
	})
	if err != nil {
		s.log.Error("seller registration failed", "email", email, "error", err)
		return seller.Seller{}, err
	}
	s.log.Info("seller registered", "seller_id", created.ID)
	return created, nil
}
func (s *SellerService) GetByID(ctx context.Context, id int64) (seller.Seller, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SellerService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Info("seller deleted", "seller_id", id)
	return nil
}

func (s *SellerService) List(ctx context.Context, f incoming.ListFilter) ([]seller.Seller, error) {
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
func (s *SellerService) Update(ctx context.Context, id int64, in incoming.UpdateSellerInput) (seller.Seller, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return seller.Seller{}, err
	}

	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return seller.Seller{}, seller.ErrInvalidInput
		}
		current.Name = name
	}

	if in.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*in.Email))
		if email == "" {
			return seller.Seller{}, seller.ErrInvalidInput
		}
		if email != strings.ToLower(current.Email) {
			_, err := s.repo.GetByEmail(ctx, email)
			if err == nil {
				return seller.Seller{}, seller.ErrEmailAlreadyExists
			}
			if !errors.Is(err, seller.ErrNotFound) {
				return seller.Seller{}, err
			}
		}
		current.Email = email
	}

	if in.Phone != nil {
		current.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Address != nil {
		current.Address = strings.TrimSpace(*in.Address)
	}

	updated, err := s.repo.Update(ctx, current)
	if err != nil {
		s.log.Error("seller update failed", "seller_id", id, "error", err)
		return seller.Seller{}, err
	}
	s.log.Info("seller updated", "seller_id", id)
	return updated, nil
}
