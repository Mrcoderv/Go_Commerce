package server

import (
	"Go_Commerce/internal/infrastructure/config"
	"context"
	"fmt"
	"log"
	nethttp "net/http"

	"go.uber.org/fx"
)

func NewServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	handler nethttp.Handler,
) *nethttp.Server {

	srv := &nethttp.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil &&
					err != nethttp.ErrServerClosed {
					log.Printf("HTTP server error: %v", err)
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
