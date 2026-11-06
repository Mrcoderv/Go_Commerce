package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type BuyerHandler struct {
	service incoming.BuyerService
}

func NewBuyerHandler(service incoming.BuyerService) *BuyerHandler {
	return &BuyerHandler{service: service}
}

func (handler *BuyerHandler) Register(ctx context.Context, r *dto.RegisterBuyerRequest) (*dto.BuyerResponse, error) {
	out, err := handler.service.Register(ctx, incoming.RegisterBuyerInput{
		Name:    r.Body.Name,
		Email:   r.Body.Email,
		Phone:   r.Body.Phone,
		Address: r.Body.Address,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (h *BuyerHandler) Get(ctx context.Context, req *dto.GetBuyerRequest) (*dto.BuyerResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	out, err := h.service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (h *BuyerHandler) List(ctx context.Context, req *dto.ListBuyersRequest) (*dto.ListBuyersResponse, error) {
	out, err := h.service.List(ctx, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(err)
	}
	if out == nil {
		out = []buyer.Buyer{}
	}
	return &dto.ListBuyersResponse{Body: out}, nil
}

func (handler *BuyerHandler) Update(ctx context.Context, req *dto.UpdateBuyerRequest) (*dto.BuyerResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	out, err := handler.service.Update(ctx, req.ID, incoming.UpdateBuyerInput{
		Name:    &req.Body.Name,
		Phone:   &req.Body.Phone,
		Address: &req.Body.Address,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (handler *BuyerHandler) Delete(ctx context.Context, r *dto.DeleteBuyerRequest) (*struct{}, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != r.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(err)
	}
	if err := handler.service.Delete(ctx, r.ID); err != nil {
		return nil, MapError(err)
	}
	return nil, nil
}
