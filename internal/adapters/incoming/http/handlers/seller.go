package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type SellerHandler struct {
	Service incoming.SellerService
}

func NewSellerHandler(s incoming.SellerService) *SellerHandler {
	return &SellerHandler{Service: s}
}

func (h *SellerHandler) Register(ctx context.Context, req *dto.RegisterSellerRequest) (*dto.SellerResponse, error) {
	out, err := h.Service.Register(ctx, incoming.RegisterSellerInput{
		Name:    req.Body.Name,
		Email:   req.Body.Email,
		Phone:   req.Body.Phone,
		Address: req.Body.Address,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.SellerResponse{Body: out}, nil
}

func (h *SellerHandler) GetByID(ctx context.Context, req *dto.GetSellerRequest) (*dto.SellerResponse, error) {
	claims, err := sellerClaims(ctx)
	if err != nil || *claims.SellerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	out, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.SellerResponse{Body: out}, nil
}

func (h *SellerHandler) List(ctx context.Context, req *dto.ListSellersRequest) (*dto.ListSellersResponse, error) {
	out, err := h.Service.List(ctx, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(err)
	}
	if out == nil {
		out = []seller.Seller{}
	}
	return &dto.ListSellersResponse{Body: out}, nil
}

func (h *SellerHandler) Update(ctx context.Context, req *dto.UpdateSellerRequest) (*dto.SellerResponse, error) {
	claims, err := sellerClaims(ctx)
	if err != nil || *claims.SellerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	out, err := h.Service.Update(ctx, req.ID, incoming.UpdateSellerInput{
		Name:    req.Body.Name,
		Email:   req.Body.Email,
		Phone:   req.Body.Phone,
		Address: req.Body.Address,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.SellerResponse{Body: out}, nil
}

func (h *SellerHandler) Delete(ctx context.Context, req *dto.DeleteSellerRequest) (*struct{}, error) {
	claims, err := sellerClaims(ctx)
	if err != nil || *claims.SellerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	if err := h.Service.Delete(ctx, req.ID); err != nil {
		return nil, MapError(err)
	}
	return nil, nil
}
