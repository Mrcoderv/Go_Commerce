package main

import (
	"Go_Commerce/internal/adapters/incoming/http"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/infrastructure/server"

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
		),
		fx.Invoke(func(srv *nethttp.Server) {}),
	)
	app.Run()
}
