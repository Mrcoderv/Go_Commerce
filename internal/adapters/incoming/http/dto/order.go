package dto

import "Go_Commerce/internal/domain/order"

type OrderItemDTO struct {
	ProductID int64 `json:"product_id" minimum:"1" doc:"Product ID"`
	Quantity  int   `json:"quantity" minimum:"1" doc:"Quantity to buy"`
}

type PlaceOrderRequest struct {
	Body struct {
		Items []OrderItemDTO `json:"items" minItems:"1" doc:"Items to order; buyer comes from the JWT"`
	}
}

type GetOrderRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type UpdateOrderStatusRequest struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Status order.Status `json:"status" enum:"pending,confirmed,shipped,delivered,cancelled,returned" doc:"Next status; must be a legal transition"`
	}
}

type CancelOrderRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type ListBuyerOrdersRequest struct {
	BuyerID int64 `path:"id" minimum:"1" doc:"Buyer ID; must match the buyer in the JWT"`
	Limit   int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset  int   `query:"offset" default:"0" minimum:"0"`
}

type ListSellerOrdersRequest struct {
	SellerID int64 `path:"id" minimum:"1" doc:"Seller ID; must match the seller in the JWT"`
	Limit    int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset   int   `query:"offset" default:"0" minimum:"0"`
}

type OrderResponse struct {
	Body order.Order
}

type ListOrdersResponse struct {
	Body []order.Order
}
