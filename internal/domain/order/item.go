package order

type Item struct {
	ID          int64  `json:"id"`
	OrderID     int64  `json:"order_id"` // was a duplicate "id" tag
	ProductID   int64  `json:"product_id"`
	SellerID    int64  `json:"seller_id"`
	ProductName string `json:"product_name"` // snapshot
	UnitPrice   int64  `json:"unit_price"`   // snapshot, paisa
	Quantity    int    `json:"quantity"`
}
