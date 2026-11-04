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
type ListProductsRequest struct {
	SellerID int64 `query:"seller_id" minimum:"0" doc:"Only this seller's products. Leave out for all."`
	Limit    int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset   int   `query:"offset" default:"0" minimum:"0"`
}

type ListProductsResponse struct {
	Body []product.Product
}

func (h *ProductHandler) List(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error) {
	f := incoming.ProductFilter{Limit: req.Limit, Offset: req.Offset}
	if req.SellerID > 0 {
		f.SellerID = &req.SellerID
	}
	out, err := h.Service.List(ctx, f)
	if err != nil {
		return nil, mapProductError(err)
	}
	if out == nil {
		out = []product.Product{}
	}
	return &ListProductsResponse{Body: out}, nil
}

type UpdateProductRequest struct {
	ID       int64 `path:"id" minimum:"1"`
	SellerID int64 `header:"X-Seller-ID" required:"true" minimum:"1" doc:"The seller who owns the product"`
	Body     struct {
		Name        *string `json:"name,omitempty" minLength:"1"`
		Description *string `json:"description,omitempty"`
		UnitPrice   *int64  `json:"unit_price,omitempty" minimum:"1" doc:"Price in paisa"`
		Stock       *int    `json:"stock,omitempty" minimum:"0"`
	}
}

type UpdateProductResponse struct {
	Body product.Product
}

func (h *ProductHandler) Update(ctx context.Context, req *UpdateProductRequest) (*UpdateProductResponse, error) {
	out, err := h.Service.Update(ctx, req.ID, req.SellerID, incoming.UpdateProductInput{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		UnitPrice:   req.Body.UnitPrice,
		Stock:       req.Body.Stock,
	})
	if err != nil {
		return nil, mapProductError(err)
	}
	return &UpdateProductResponse{Body: out}, nil
}

type DeleteProductRequest struct {
	ID       int64 `path:"id" minimum:"1"`
	SellerID int64 `header:"X-Seller-ID" required:"true" minimum:"1" doc:"The seller who owns the product"`
}

type DeleteProductResponse struct{}

func (h *ProductHandler) Delete(ctx context.Context, req *DeleteProductRequest) (*DeleteProductResponse, error) {
	if err := h.Service.Delete(ctx, req.ID, req.SellerID); err != nil {
		return nil, mapProductError(err)
	}
	return &DeleteProductResponse{}, nil
}