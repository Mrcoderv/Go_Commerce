package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type BuyerHandler struct {
	service incoming.BuyerService
	log     logger.Logger
}

func NewBuyerHandler(service incoming.BuyerService, log logger.Logger) *BuyerHandler {
	return &BuyerHandler{service: service, log: log}
}

func (h *BuyerHandler) Register(ctx context.Context, r *dto.RegisterBuyerRequest) (*dto.BuyerResponse, error) {
	out, err := h.service.Register(ctx, incoming.RegisterBuyerInput{
		Name:    r.Body.Name,
		Email:   r.Body.Email,
		Phone:   r.Body.Phone,
		Address: r.Body.Address,
	})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (h *BuyerHandler) Get(ctx context.Context, req *dto.GetBuyerRequest) (*dto.BuyerResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(h.log, err)
	}
	out, err := h.service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (h *BuyerHandler) List(ctx context.Context, req *dto.ListBuyersRequest) (*dto.ListBuyersResponse, error) {
	out, err := h.service.List(ctx, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if out == nil {
		out = []buyer.Buyer{}
	}
	return &dto.ListBuyersResponse{Body: out}, nil
}

func (h *BuyerHandler) Update(ctx context.Context, req *dto.UpdateBuyerRequest) (*dto.BuyerResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != req.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(h.log, err)
	}
	out, err := h.service.Update(ctx, req.ID, incoming.UpdateBuyerInput{
		Name:    req.Body.Name,
		Phone:   req.Body.Phone,
		Address: req.Body.Address,
	})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.BuyerResponse{Body: out}, nil
}

func (h *BuyerHandler) Delete(ctx context.Context, r *dto.DeleteBuyerRequest) (*struct{}, error) {
	claims, err := buyerClaims(ctx)
	if err != nil || *claims.BuyerID != r.ID {
		if err == nil {
			err = domainauth.ErrForbidden
		}
		return nil, MapError(h.log, err)
	}
	if err := h.service.Delete(ctx, r.ID); err != nil {
		return nil, MapError(h.log, err)
	}
	return nil, nil
}
