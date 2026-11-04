package handlers

import (
	"context"
	"errors"

	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/ports/incoming"

	"github.com/danielgtaylor/huma/v2"
)

type BuyerHandler struct {
	service incoming.BuyerService
}

func NewBuyerHandler(service incoming.BuyerService) *BuyerHandler {
	return &BuyerHandler{service: service}
}

// json request body to the browser
type RegisterBuyerRequest struct {
	Body struct {
		Name    string `json:"name" minlength:"1"`
		Email   string `json:"email" format:"email"`
		Phone   string `json:"phone,omitempty"` // omitempty is used hide the valete os there is no data.
		Address string `json:"address"`
	}
}
type BuyerResponse struct {
	Body buyer.Buyer
}

func (handler *BuyerHandler) Register(ctx context.Context, r *RegisterBuyerRequest) (*BuyerResponse, error) {
	outing, err := handler.service.Register(ctx, incoming.RegisterBuyerInput{
		Name:    r.Body.Name,
		Email:   r.Body.Email,
		Phone:   r.Body.Phone,
		Address: r.Body.Address,
	})
	if err != nil {
		return nil, mapBuyerError(err)
	}

	return &BuyerResponse{Body: outing}, nil
}

type GetBuyerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

func (h *BuyerHandler) Get(ctx context.Context, req *GetBuyerRequest) (*BuyerResponse, error) {
	out, err := h.service.GetByID(ctx, req.ID)

	if err != nil {
		return nil, mapSellerError(err)

	}
	return &BuyerResponse{Body: out}, nil
}

// updtte
type UpdateBuyerRequest struct {
	ID   int64 `path:"id"`
	Body struct {
		Name    string `json:"name,omitempty"`
		Phone   string `json:"phone,omitempty"` // omitempty is used hide the valete os there is no data.
		Address string `json:"address,omitempty"`
	}
}

func (handler *BuyerHandler) Update(ctx context.Context, req *UpdateBuyerRequest) (*BuyerResponse, error) {
	out, err := handler.service.Update(ctx, req.ID, incoming.UpdateBuyerInput{
		Name:    &req.Body.Name,
		Phone:   &req.Body.Phone,
		Address: &req.Body.Address, // & passing the service layer as the pointer.
	})
	if err != nil {
		return nil, mapBuyerError(err)
	}
	return &BuyerResponse{Body: out}, nil
}

// delete
type DeleteBuyerRequest struct {
	ID int64 `path:"id"`
}

func (handler *BuyerHandler) Delete(ctx context.Context, r *DeleteBuyerRequest) (*struct{}, error) {
	if err := handler.service.Delete(ctx, r.ID); err != nil {
		return nil, mapSellerError(err)
	}
	return nil, nil
}

// maping the error to the doiam error into the http error
func mapBuyerError(err error) error {
	switch {
	case errors.Is(err, buyer.ErrNotFound):
		return huma.Error404NotFound("buyer not found")
	case errors.Is(err, buyer.ErrEmailAlreadyExists):
		return huma.Error409Conflict("email already exists")
	case errors.Is(err, buyer.ErrInvalidInput):
		return huma.Error400BadRequest("invalid buyer input")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}
