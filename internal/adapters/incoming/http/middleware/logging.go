package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		huma.LogRequest(r)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
