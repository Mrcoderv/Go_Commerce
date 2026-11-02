package seller

type Seller struct {
	SellerID int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    int64  `json:"phone"`
	Address  string `json:"address"`
}
