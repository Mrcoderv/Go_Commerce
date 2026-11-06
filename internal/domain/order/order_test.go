package order

import "testing"

func TestForSellerHidesOtherSellersItems(t *testing.T) {
	o := Order{
		ID:      1,
		BuyerID: 2,
		Total:   4000,
		Items: []Item{
			{ProductID: 10, SellerID: 5, UnitPrice: 1000, Quantity: 2}, // 2000
			{ProductID: 20, SellerID: 7, UnitPrice: 2000, Quantity: 1}, // 2000
		},
	}

	got := ForSeller(o, 5)
	if len(got.Items) != 1 || got.Items[0].SellerID != 5 {
		t.Fatalf("expected only seller 5's item, got %+v", got.Items)
	}
	if got.Total != 2000 {
		t.Errorf("expected seller 5 subtotal 2000, got %d", got.Total)
	}

	// The original order must not be mutated.
	if len(o.Items) != 2 || o.Total != 4000 {
		t.Errorf("ForSeller mutated the original order: %+v", o)
	}
}

func TestForSellerWithoutItems(t *testing.T) {
	o := Order{ID: 1, Total: 4000, Items: []Item{{ProductID: 10, SellerID: 5, UnitPrice: 4000, Quantity: 1}}}

	got := ForSeller(o, 99)
	if len(got.Items) != 0 {
		t.Fatalf("expected no items, got %+v", got.Items)
	}
	if got.Total != 0 {
		t.Errorf("expected zero subtotal, got %d", got.Total)
	}
}
