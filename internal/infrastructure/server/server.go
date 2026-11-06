package server

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"

	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"

	"go.uber.org/fx"
)

func NewServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	handler nethttp.Handler,
	log logger.Logger,
) *nethttp.Server {

	srv := &nethttp.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("starting http server", "addr", srv.Addr, "env", cfg.AppEnv)
			go func() {
				if err := srv.ListenAndServe(); err != nil &&
					!errors.Is(err, nethttp.ErrServerClosed) {
					log.Error("http server error", "error", err)
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			log.Info("shutting down http server", "addr", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
