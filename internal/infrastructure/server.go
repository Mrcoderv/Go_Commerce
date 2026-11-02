package server

import (
	"Go_Commerce/internal/infrastructure/config"
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"
)

func NewServer(
	lc fx.Lifecycle,
	c config.Config,
) *http.Server {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.AppPort),
		Handler: http.DefaultServeMux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil &&
					err != http.ErrServerClosed {
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
