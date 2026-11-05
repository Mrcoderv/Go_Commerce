package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	"Go_Commerce/internal/adapters/incoming/http/middleware"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type ProductHandler struct {
	Service incoming.ProductService
}

func NewProductHandler(s incoming.ProductService) *ProductHandler {
	return &ProductHandler{Service: s}
}

func sellerIDFromAuth(ctx context.Context) (int64, bool) {
	c, ok := middleware.ClaimsFromContext(ctx)
	if !ok || c.Role != domainauth.RoleSeller || c.SellerID == nil {
		return 0, false
	}
	return *c.SellerID, true
}

func (h *ProductHandler) Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	sellerID := req.Body.SellerID
	if sellerID == 0 {
		sellerID = req.SellerIDHeader
	}
	if id, ok := sellerIDFromAuth(ctx); ok {
		sellerID = id
	}
	if sellerID == 0 {
		return nil, MapError(domainauth.ErrUnauthorized)
	}

	out, err := h.Service.Create(ctx, incoming.CreateProductInput{
		SellerID:    sellerID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Price:       req.Body.Price,
		Stock:       req.Body.Stock,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.ProductResponse{Body: out}, nil
}

func (h *ProductHandler) GetByID(ctx context.Context, req *dto.GetProductRequest) (*dto.ProductResponse, error) {
	out, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.ProductResponse{Body: out}, nil
}

func (h *ProductHandler) List(ctx context.Context, req *dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
	var sellerID *int64
	if req.SellerID > 0 {
		sellerID = &req.SellerID
	}
	out, err := h.Service.List(ctx, incoming.ProductFilter{
		SellerID: sellerID,
		Limit:    req.Limit,
		Offset:   req.Offset,
	})
	if err != nil {
		return nil, MapError(err)
	}
	if out == nil {
		out = []product.Product{}
	}
	return &dto.ListProductsResponse{Body: out}, nil
}

func (h *ProductHandler) Update(ctx context.Context, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	sellerID := req.Body.SellerID
	if sellerID == 0 {
		sellerID = req.SellerIDHeader
	}
	if id, ok := sellerIDFromAuth(ctx); ok {
		sellerID = id
	}
	if sellerID == 0 {
		return nil, MapError(domainauth.ErrUnauthorized)
	}

	var priceInPaisa *float64
	if req.Body.Price != nil {
		p := *req.Body.Price
		priceInPaisa = &p
	}

	out, err := h.Service.Update(ctx, req.ID, sellerID, incoming.UpdateProductInput{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		Price:       priceInPaisa,
		Stock:       req.Body.Stock,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.ProductResponse{Body: out}, nil
}

func (h *ProductHandler) Delete(ctx context.Context, req *dto.DeleteProductRequest) (*struct{}, error) {
	sellerID := req.SellerIDHeader
	if sellerID == 0 {
		sellerID = req.SellerIDQuery
	}
	if id, ok := sellerIDFromAuth(ctx); ok {
		sellerID = id
	}
	if sellerID == 0 {
		return nil, MapError(domainauth.ErrUnauthorized)
	}
	if err := h.Service.Delete(ctx, req.ID, sellerID); err != nil {
		return nil, MapError(err)
	}
	return nil, nil
}
