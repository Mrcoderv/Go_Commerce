package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type OrderHandler struct {
	Service incoming.OrderService
	log     logger.Logger
}

func NewOrderHandler(service incoming.OrderService, log logger.Logger) *OrderHandler {
	return &OrderHandler{Service: service, log: log}
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, req *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	items := make([]incoming.OrderItemInput, len(req.Body.Items))
	for i, it := range req.Body.Items {
		items[i] = incoming.OrderItemInput{ProductID: it.ProductID, Quantity: it.Quantity}
	}

	o, err := h.Service.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: *claims.BuyerID,
		Items:   items,
	})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *dto.GetOrderRequest) (*dto.OrderResponse, error) {
	claims, err := orderAccessClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	o, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if !canAccessOrder(claims, o) {
		return nil, MapError(h.log, domainauth.ErrForbidden)
	}
	if claims.Role == domainauth.RoleSeller {
		o = order.ForSeller(o, *claims.SellerID)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) UpdateStatus(ctx context.Context, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	claims, err := sellerClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	current, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if !canAccessOrder(claims, current) {
		return nil, MapError(h.log, domainauth.ErrForbidden)
	}
	o, err := h.Service.UpdateStatus(ctx, req.ID, req.Body.Status)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.OrderResponse{Body: order.ForSeller(o, *claims.SellerID)}, nil
}

func (h *OrderHandler) Cancel(ctx context.Context, req *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	current, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if !canAccessOrder(claims, current) {
		return nil, MapError(h.log, domainauth.ErrForbidden)
	}
	if err := h.Service.Cancel(ctx, req.ID); err != nil {
		return nil, MapError(h.log, err)
	}
	o, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) ListByBuyer(ctx context.Context, req *dto.ListBuyerOrdersRequest) (*dto.ListOrdersResponse, error) {
	claims, err := buyerClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if req.BuyerID != *claims.BuyerID {
		return nil, MapError(h.log, domainauth.ErrForbidden)
	}
	orders, err := h.Service.ListByBuyer(ctx, *claims.BuyerID, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if orders == nil {
		orders = []order.Order{}
	}
	return &dto.ListOrdersResponse{Body: orders}, nil
}

func (h *OrderHandler) ListBySeller(ctx context.Context, req *dto.ListSellerOrdersRequest) (*dto.ListOrdersResponse, error) {
	claims, err := sellerClaims(ctx)
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if req.SellerID != *claims.SellerID {
		return nil, MapError(h.log, domainauth.ErrForbidden)
	}
	orders, err := h.Service.ListBySeller(ctx, *claims.SellerID, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	if orders == nil {
		orders = []order.Order{}
	}
	return &dto.ListOrdersResponse{Body: orders}, nil
}
