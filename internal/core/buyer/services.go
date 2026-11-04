package buyer

import (
	"context"
	"errors"
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
func (s *BuyerService) Register(
    ctx context.Context,
    in incoming.RegisterBuyerInput,
) (buyer.Buyer, error) {
	// spacig remove
    name := strings.TrimSpace(in.Name)
    email := strings.ToLower(strings.TrimSpace(in.Email))

    if name == "" || email == "" {
        return buyer.Buyer{}, buyer.ErrInvalidInput
    }

    // Check whether the email is valid
    if _, err := mail.ParseAddress(email); err != nil {
        return buyer.Buyer{}, buyer.ErrInvalidInput
    }

    // Check if a buyer with this email already exists
    existingBuyer, err := s.repo.GetByEmail(ctx, email)
    if err == nil && existingBuyer.ID != 0 {
        return buyer.Buyer{}, buyer.ErrEmailAlreadyExists
    }

    if err != nil && !errors.Is(err, buyer.ErrNotFound) {
        return buyer.Buyer{}, err
    }

    // Create the new buyer
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
