package tests

import (
	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/adapters/outgoing/postgres/repositories"
	"Go_Commerce/internal/core/buyer"
	"Go_Commerce/internal/core/order"
	"Go_Commerce/internal/core/product"
	"Go_Commerce/internal/core/seller"
	orderdomain "Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/migrations"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func TestOrderFlow_Integration(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("skipping integration test: TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	// Apply migrations (safe to run multiple times)
	if err := migrations.Run(ctx, sqlDB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Setup adapters and core services
	txManager := postgres.NewTxManager(sqlDB)
	sellerRepo := repositories.NewSellerRepository(sqlDB)
	productRepo := repositories.NewProductRepository(sqlDB)
	buyerRepo := repositories.NewBuyerRepository(sqlDB)
	orderRepo := repositories.NewOrderRepository(sqlDB)

	sellerSvc := seller.NewSellerService(sellerRepo)
	productSvc := product.NewProductService(productRepo)
	buyerSvc := buyer.NewBuyerService(buyerRepo)
	orderSvc := order.NewOrderService(orderRepo, buyerRepo, productRepo, txManager)

	now := time.Now().UnixNano()

	// 1. Create Seller with unique email
	s, err := sellerSvc.Register(ctx, incoming.RegisterSellerInput{
		Name:    "Integration Seller",
		Email:   fmt.Sprintf("seller-%d@example.com", now),
		Phone:   "+9779800000001",
		Address: "Kathmandu",
	})
	if err != nil {
		t.Fatalf("failed to register seller: %v", err)
	}

	// 2. Create Product (price in paisa: 50000 = NPR 500, stock: 10)
	p, err := productSvc.Create(ctx, incoming.CreateProductInput{
		SellerID:    s.ID,
		Name:        "Test Product",
		Description: "Integration test product description",
		Price:       50000,
		Stock:       10,
	})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// 3. Create Buyer with unique email
	b, err := buyerSvc.Register(ctx, incoming.RegisterBuyerInput{
		Name:    "Integration Buyer",
		Email:   fmt.Sprintf("buyer-%d@example.com", now),
		Phone:   "+9779800000002",
		Address: "Lalitpur",
	})
	if err != nil {
		t.Fatalf("failed to register buyer: %v", err)
	}

	// 4. Place Order for 3 units
	o, err := orderSvc.PlaceOrder(ctx, incoming.PlaceOrderInput{
		BuyerID: b.ID,
		Items: []incoming.OrderItemInput{
			{ProductID: p.ID, Quantity: 3},
		},
	})
	if err != nil {
		t.Fatalf("failed to place order: %v", err)
	}
	if o.Total != 150000 {
		t.Errorf("expected total 150000 paisa, got %d", o.Total)
	}

	// 5. Check stock dropped from 10 to 7
	pAfterOrder, err := productSvc.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}
	if pAfterOrder.Stock != 7 {
		t.Errorf("expected stock to drop to 7, got %d", pAfterOrder.Stock)
	}

	// 6. Cancel the order
	err = orderSvc.Cancel(ctx, o.ID)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	// 7. Verify order status cancelled
	oCancelled, err := orderSvc.GetByID(ctx, o.ID)
	if err != nil {
		t.Fatalf("failed to get cancelled order: %v", err)
	}
	if oCancelled.Status != orderdomain.StatusCancelled {
		t.Errorf("expected status cancelled, got %s", oCancelled.Status)
	}

	// 8. Check stock restored back to 10
	pAfterCancel, err := productSvc.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("failed to get product after cancel: %v", err)
	}
	if pAfterCancel.Stock != 10 {
		t.Errorf("expected stock to be restored to 10, got %d", pAfterCancel.Stock)
	}
}
