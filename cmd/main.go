package main

import (
	nethttp "net/http"

	"Go_Commerce/internal/adapters/incoming/http"
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/adapters/outgoing/postgres/repositories"
	"Go_Commerce/internal/core/buyer"
	"Go_Commerce/internal/core/order"
	"Go_Commerce/internal/core/product"
	"Go_Commerce/internal/core/seller"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/database"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/infrastructure/server"
	"Go_Commerce/internal/ports/outgoing"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			logger.NewLogger,
			http.New,
			server.NewServer,
			postgres.NewTxManager,

			// repositories, each bound to its port interface
			fx.Annotate(repositories.NewSellerRepository, fx.As(new(outgoing.SellerRepository))),
			fx.Annotate(repositories.NewProductRepository, fx.As(new(outgoing.ProductRepository))),
			fx.Annotate(repositories.NewBuyerRepository, fx.As(new(outgoing.BuyerRepository))),
			fx.Annotate(repositories.NewOrderRepository, fx.As(new(outgoing.OrderRepository))),

			// handlers
			handlers.NewSellerHandler,
			handlers.NewProductHandler,
			handlers.NewBuyerHandler,
			handlers.NewOrderHandler,
		),
		database.Module,
		seller.Module,
		product.Module,
		buyer.Module,
		order.Module,
		fx.Invoke(func(srv *nethttp.Server) {}),
	)
	app.Run()
}