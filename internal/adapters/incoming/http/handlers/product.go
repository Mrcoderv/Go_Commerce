package handlers

import (
	"context"
	"errors"

	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"

	"github.com/danielgtaylor/huma/v2"
)

type ProductHandler struct {
	Service incoming.ProductService
}

func NewProductHandler(s incoming.ProductService) *ProductHandler {
	return &ProductHandler{Service: s}
}

type CreateProductRequest struct {
	Body struct {
		SellerID    int64  `json:"seller_id" minimum:"1"`
		Name        string `json:"name" minLength:"1"`
		Description string `json:"description,omitempty"`
		UnitPrice   int64  `json:"unit_price" minimum:"1" doc:"Price in paisa"`
		Stock       int    `json:"stock" minimum:"0"`
	}
}

type CreateProductResponse struct {
	Body product.Product
}

func (h *ProductHandler) Create(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	out, err := h.Service.Create(ctx, incoming.CreateProductInput{
		SellerID:    req.Body.SellerID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		UnitPrice:   req.Body.UnitPrice,
		Stock:       req.Body.Stock,
	})
	if err != nil {
		return nil, mapProductError(err)
	}
	return &CreateProductResponse{Body: out}, nil
}

type GetProductRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type GetProductResponse struct {
	Body product.Product
}

func (h *ProductHandler) GetByID(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	out, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapProductError(err)
	}
	return &GetProductResponse{Body: out}, nil
}

func mapProductError(err error) error {
	switch {
	case errors.Is(err, product.ErrNotFound):
		return huma.Error404NotFound("product not found")
	case errors.Is(err, product.ErrInvalidInput):
		return huma.Error400BadRequest("invalid input")
	case errors.Is(err, product.ErrForbidden):
		return huma.Error403Forbidden("you do not own this product")
	case errors.Is(err, product.ErrInsufficientStock):
		return huma.Error409Conflict("insufficient stock")
	default:
		return huma.Error500InternalServerError("internal error")
	}
}