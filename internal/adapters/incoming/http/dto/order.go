package dto

import "Go_Commerce/internal/domain/order"

type OrderItemDTO struct {
	ProductID int64 `json:"product_id" minimum:"1"`
	Quantity  int   `json:"quantity" minimum:"1"`
}

type PlaceOrderRequest struct {
	Body struct {
		BuyerID int64          `json:"buyer_id,omitempty"`
		Items   []OrderItemDTO `json:"items" minItems:"1"`
	}
}

type GetOrderRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type UpdateOrderStatusRequest struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Status order.Status `json:"status" enum:"pending,confirmed,shipped,delivered,cancelled"`
	}
}

type CancelOrderRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type ListBuyerOrdersRequest struct {
	BuyerID int64 `path:"id" minimum:"1"`
	Limit   int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset  int   `query:"offset" default:"0" minimum:"0"`
}

type ListSellerOrdersRequest struct {
	SellerID int64 `path:"id" minimum:"1"`
	Limit    int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset   int   `query:"offset" default:"0" minimum:"0"`
}

type OrderResponse struct {
	Body order.Order
}

type ListOrdersResponse struct {
	Body []order.Order
}
