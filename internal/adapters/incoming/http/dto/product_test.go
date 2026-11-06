package dto

import (
	"encoding/json"
	"testing"
)

func TestCreateProductRequest_PriceIsIntegerPaisa(t *testing.T) {
	var req CreateProductRequest
	err := json.Unmarshal([]byte(`{"Body":{"name":"Laptop","price":49950,"stock":1}}`), &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Body.Price != 49950 {
		t.Errorf("expected 49950 paisa, got %d", req.Body.Price)
	}
}

func TestCreateProductRequest_RejectsFractionalPaisa(t *testing.T) {
	var req CreateProductRequest
	err := json.Unmarshal([]byte(`{"Body":{"name":"Laptop","price":499.50,"stock":1}}`), &req)
	if err == nil {
		t.Fatal("expected fractional paisa to be rejected, got nil error")
	}
}

func TestUpdateProductRequest_RejectsFractionalPaisa(t *testing.T) {
	var req UpdateProductRequest
	err := json.Unmarshal([]byte(`{"Body":{"price":499.50}}`), &req)   
	if err == nil {
		t.Fatal("expected fractional paisa to be rejected, got nil error")
	}
}
