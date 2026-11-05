package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	"Go_Commerce/internal/adapters/incoming/http/middleware"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type OrderHandler struct {
	Service incoming.OrderService
}

func NewOrderHandler(service incoming.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func buyerIDFromAuth(ctx context.Context) (int64, bool) {
	c, ok := middleware.ClaimsFromContext(ctx)
	if !ok || c.Role != domainauth.RoleBuyer || c.BuyerID == nil {
		return 0, false
	}
	return *c.BuyerID, true
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, req *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	buyerID := req.Body.BuyerID
	if id, ok := buyerIDFromAuth(ctx); ok {
		buyerID = id
	}
	if buyerID == 0 {
		return nil, MapError(domainauth.ErrUnauthorized)
	}
	items := make([]incoming.OrderItemInput, len(req.Body.Items))
	for i, it := range req.Body.Items {
		items[i] = incoming.OrderItemInput{ProductID: it.ProductID, Quantity: it.Quantity}
	}

	o, err := h.Service.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: buyerID,
		Items:   items,
	})
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *dto.GetOrderRequest) (*dto.OrderResponse, error) {
	o, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) UpdateStatus(ctx context.Context, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	o, err := h.Service.UpdateStatus(ctx, req.ID, req.Body.Status)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) Cancel(ctx context.Context, req *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	if err := h.Service.Cancel(ctx, req.ID); err != nil {
		return nil, MapError(err)
	}
	o, err := h.Service.GetByID(ctx, req.ID)
	if err != nil {
		return nil, MapError(err)
	}
	return &dto.OrderResponse{Body: o}, nil
}

func (h *OrderHandler) ListByBuyer(ctx context.Context, req *dto.ListBuyerOrdersRequest) (*dto.ListOrdersResponse, error) {
	orders, err := h.Service.ListByBuyer(ctx, req.BuyerID, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(err)
	}
	if orders == nil {
		orders = []order.Order{}
	}
	return &dto.ListOrdersResponse{Body: orders}, nil
}

func (h *OrderHandler) ListBySeller(ctx context.Context, req *dto.ListSellerOrdersRequest) (*dto.ListOrdersResponse, error) {
	orders, err := h.Service.ListBySeller(ctx, req.SellerID, incoming.ListFilter{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		return nil, MapError(err)
	}
	if orders == nil {
		orders = []order.Order{}
	}
	return &dto.ListOrdersResponse{Body: orders}, nil
}
