package handlers

import (
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/incoming"
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

type SellerHandler struct {
	Service incoming.SellerService
}
type RegisterSellerRequest struct {
	Body struct {
		Name    string `json:"name" minLength:"1"`
		Email   string `json:"email" format:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}
}

type RegisterSellerResponse struct {
	Body seller.Seller
}

func (h *SellerHandler) Register(ctx context.Context, req *RegisterSellerRequest) (*RegisterSellerResponse, error) {
	out, err := h.Service.Register(ctx, incoming.RegisterSellerInput{
		Name:    req.Body.Name,
		Email:   req.Body.Email,
		Phone:   req.Body.Phone,
		Address: req.Body.Address,
	})
	if err != nil {
		return nil, mapSellerError(err)
	}
	return &RegisterSellerResponse{Body: out}, nil
}
func mapSellerError(err error) error {
	switch {
	case errors.Is(err, seller.ErrNotFound):
		return huma.Error404NotFound("seller not found")
	case errors.Is(err, seller.ErrEmailAlreadyExists):
		return huma.Error409Conflict("email already exists")
	case errors.Is(err, seller.ErrInvalidInput):
		return huma.Error400BadRequest("invalid input")
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
