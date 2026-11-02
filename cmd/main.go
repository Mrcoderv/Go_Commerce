package main

import (
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/infrastructure/server"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			logger.NewLogger,
			server.NewServer,
		),
	)
	app.Run()
}