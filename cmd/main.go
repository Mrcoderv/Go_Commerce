package main

import (
	"Go_Commerce/internal/adapters/incoming/http"
	"Go_Commerce/internal/adapters/outgoing/memory"
	"Go_Commerce/internal/core/seller"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/database"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/infrastructure/server"
	"fmt"
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
			seller.NewSellerService,
			// NewStubSellerRepo,
			memory.NewSellerRepo,
		),
		database.Module,
		fx.Invoke(func(db *database.Postgres) {
			fmt.Println("Database connection established:", db.DB)
		}),
		fx.Invoke(func(srv *nethttp.Server) {}),
	)
	app.Run()
}
