package main

import (
	"Go_Commerce/internal/adapters/incoming/http"
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/adapters/outgoing/memory"
	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/adapters/outgoing/postgres/repositories"
	"Go_Commerce/internal/core/buyer"
	"Go_Commerce/internal/core/order"
	"Go_Commerce/internal/core/seller"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/database"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/infrastructure/server"
	"Go_Commerce/internal/ports/outgoing"
	nethttp "net/http"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			logger.NewLogger,
			http.New,
			server.NewServer,

			// Developer A's part (unchanged for now)
			seller.NewSellerService,
			memory.NewSellerRepo,

			// Buyer and order: repositories (bound to their port interfaces)
			fx.Annotate(repositories.NewBuyerRepository, fx.As(new(outgoing.BuyerRepository))),
			fx.Annotate(repositories.NewOrderRepository, fx.As(new(outgoing.OrderRepository))),
			fx.Annotate(repositories.NewProductRepository, fx.As(new(outgoing.ProductRepository))),
			postgres.NewTxManager,

			// Buyer and order: handlers
			handlers.NewBuyerHandler,
			handlers.NewOrderHandler,
		),
		database.Module,
		buyer.Module, // incomming buyer modue
		order.Module, // incomming order module

		fx.Invoke(func(srv *nethttp.Server) {}),
	)
	app.Run()
}
