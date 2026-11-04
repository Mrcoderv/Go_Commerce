package handlers

import (
	"context"
	"errors"

	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"

	"github.com/danielgtaylor/huma/v2"
)

type OrderHandler struct {
	Service incoming.OrderService
}

func NewOrderHandler(service incoming.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

type GetOrderInput struct {
	ID int64 `path:"id" minimum:"1"`
}

type OrderResponse struct {
	Body order.Order
}

func (h *OrderHandler) GetOrder(ctx context.Context, input *GetOrderInput) (*OrderResponse, error) {
	o, err := h.Service.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapOrderError(err)
	}
	return &OrderResponse{Body: o}, nil
}

type ListBuyerOrdersInput struct {
	BuyerID int64 `path:"id" minimum:"1"`
	Limit   int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset  int   `query:"offset" default:"0" minimum:"0"`
}

type OrderListResponse struct {
	Body []order.Order
}

func (h *OrderHandler) ListByBuyer(ctx context.Context, in *ListBuyerOrdersInput) (*OrderListResponse, error) {
	orders, err := h.Service.ListByBuyer(ctx, in.BuyerID, incoming.ListFilter{
		Limit:  in.Limit,
		Offset: in.Offset,
	})
	if err != nil {
		return nil, mapOrderError(err)
	}
	if orders == nil {
		orders = []order.Order{} // return [] instead of null
	}
	return &OrderListResponse{Body: orders}, nil
}

type UpdateOrderStatusInput struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Status order.Status `json:"status" enum:"pending,confirmed,shipped,delivered,cancelled"`
	}
}

func (h *OrderHandler) UpdateStatus(ctx context.Context, in *UpdateOrderStatusInput) (*OrderResponse, error) {
	o, err := h.Service.UpdateStatus(ctx, in.ID, in.Body.Status)
	if err != nil {
		return nil, mapOrderError(err)
	}
	return &OrderResponse{Body: o}, nil
}

// mapOrderError turns domain errors into HTTP errors.
func mapOrderError(err error) error {
	switch {
	case errors.Is(err, order.ErrNotFound):
		return huma.Error404NotFound("order not found")
	case errors.Is(err, order.ErrInvalidInput):
		return huma.Error400BadRequest("invalid order input")
	case errors.Is(err, order.ErrInvalidTransaction):
		return huma.Error409Conflict("invalid status transition")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}
