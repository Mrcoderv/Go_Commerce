package order

import (
	"time"
)

type Order struct {
	Order_ID   int64     `json:"id"`
	SellerID   int64     `json:"seller_id"`
	buyerID    int64     `json:"buyer_id"`
	order_date time.Time `json:"order_date"`
	amount     float64   `json:"amount"`
}
