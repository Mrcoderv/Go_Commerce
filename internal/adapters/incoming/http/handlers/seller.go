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
func NewSellerHandler(s incoming.SellerService) *SellerHandler {
	return &SellerHandler{Service: s}
}
type GetSellerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type GetSellerResponse struct {
	Body seller.Seller
}

func (h *SellerHandler) GetByID(ctx context.Context, req *GetSellerRequest) (*GetSellerResponse, error) {
	out, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapSellerError(err)
	}
	return &GetSellerResponse{Body: out}, nil
}
type ListSellersRequest struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset int `query:"offset" default:"0" minimum:"0"`
}

type ListSellersResponse struct {
	Body []seller.Seller
}

func (h *SellerHandler) List(ctx context.Context, req *ListSellersRequest) (*ListSellersResponse, error) {
	out, err := h.Service.List(ctx, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, mapSellerError(err)
	}
	if out == nil {
		out = []seller.Seller{}
	}
	return &ListSellersResponse{Body: out}, nil
}
type UpdateSellerRequest struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Name    *string `json:"name,omitempty" minLength:"1"`
		Email   *string `json:"email,omitempty" format:"email"`
		Phone   *string `json:"phone,omitempty"`
		Address *string `json:"address,omitempty"`
	}
}

type UpdateSellerResponse struct {
	Body seller.Seller
}

func (h *SellerHandler) Update(ctx context.Context, req *UpdateSellerRequest) (*UpdateSellerResponse, error) {
	out, err := h.Service.Update(ctx, req.ID, incoming.UpdateSellerInput{
		Name:    req.Body.Name,
		Email:   req.Body.Email,
		Phone:   req.Body.Phone,
		Address: req.Body.Address,
	})
	if err != nil {
		return nil, mapSellerError(err)
	}
	return &UpdateSellerResponse{Body: out}, nil
}

type DeleteSellerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type DeleteSellerResponse struct{}

func (h *SellerHandler) Delete(ctx context.Context, req *DeleteSellerRequest) (*DeleteSellerResponse, error) {
	if err := h.Service.Delete(ctx, req.ID); err != nil {
		return nil, mapSellerError(err)
	}
	return &DeleteSellerResponse{}, nil
}