package order

type Order_item struct {
	Order_item_ID int64   `json:"id"`
	Order_ID      int64   `json:"id"`
	Product_ID    int64   `json:"product_id"`
	Quantity      int64   `json:"quantity"`
	Price         float64 `json:"price"`
	BuyerID       int64   `json:"buyer_id"`
}
