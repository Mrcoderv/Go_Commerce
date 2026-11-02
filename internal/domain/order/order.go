package order

import "time"

type Order struct {
	ID        int64     `json:"id"`
	BuyerID   int64     `json:"buyer_id"`
	Status    Status    `json:"status"`
	Total     int64     `json:"total"` // paisa
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}
