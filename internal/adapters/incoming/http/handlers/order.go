package handlers

import (
	"context"

	"Go_Commerce/internal/core/order"

	"github.com/danielgtaylor/huma/v2"
)

type OrderHandler struct {
	Service *order.OrderService
}

type GetOrderInput struct {
	ID int64 `path:"id"`
}

type GetOrderOutput struct {
	Body incomingOrderBody
}

type incomingOrderBody struct {
	Order any `json:"order"`
}

func (h *OrderHandler) GetOrder(ctx context.Context, input *GetOrderInput) (*GetOrderOutput, error) {
	o, err := h.Service.GetByID(ctx, input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("order not found")
	}
	return &GetOrderOutput{Body: incomingOrderBody{Order: o}}, nil
}

//
//func (h *OrderHandler) PostOrder(ctx context.Context, input *PostOrderInput) (*PostOrderOutput, error) {
//	o, err := h.Service.PostByID(ctx, input.ID)
//	if err != nil {
//		return nil, huma.Error404NotFound("order not found")
//	}
//	return &PostOrderOutput{Body: incomingOrderBody{Order: o}}, nil
//}
