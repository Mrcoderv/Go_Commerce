package order

import "time"

type Order struct {
	ID        int64     `json:"id"`
	BuyerID   int64     `json:"buyer_id"`
	Total     int64     `json:"total"` // total amount in paisaaaaa
	Status    Status    `json:"status"`
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}
