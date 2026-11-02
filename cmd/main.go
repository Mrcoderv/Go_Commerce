package main

import (
	"Go_Commerce/internal/infrastructure/server"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			server.NewServer,

		),
	)

	app.Run()
}
