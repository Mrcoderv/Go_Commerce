package product

type Product struct {
	ID          int64  `json:"id"`
	SellerID    int64  `json:"seller_id"`
	Name        string `json:"product_name"`
	Description string `json:"product_description"`
	UnitPrice   int64  `json:"unit_price"`
	Stock       int    `json:"stock"`
}
