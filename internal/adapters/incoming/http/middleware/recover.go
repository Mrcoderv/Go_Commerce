package middleware

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Info("Recovered from", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				
			}
		}()
		next.ServeHTTP(w, r)
	})
}
