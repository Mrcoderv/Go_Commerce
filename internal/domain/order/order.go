package order

import "time"

type Order struct {
	ID        int64     `json:"id"`
	BuyerID   int64     `json:"buyer_id"`
	Total     int64     `json:"total"` // total amount in paisa
	Status    Status    `json:"status"`
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}

// ForSeller returns a copy of the order that only contains items belonging to
// sellerID and a total recomputed from those items. Sellers must never see
// other sellers' items or the full order total.
func ForSeller(o Order, sellerID int64) Order {
	items := make([]Item, 0, len(o.Items))
	var total int64
	for _, it := range o.Items {
		if it.SellerID != sellerID {
			continue
		}
		items = append(items, it)
		total += it.UnitPrice * int64(it.Quantity)
	}
	o.Items = items
	o.Total = total
	return o
}
