package middleware

import (
	"context"
	"net/http"
	"strings"

	coreauth "Go_Commerce/internal/core/auth"
	"Go_Commerce/internal/infrastructure/config"
)

type authKey struct{}

func WithClaims(ctx context.Context, c *coreauth.Claims) context.Context {
	return context.WithValue(ctx, authKey{}, c)
}

func ClaimsFromContext(ctx context.Context) (*coreauth.Claims, bool) {
	c, ok := ctx.Value(authKey{}).(*coreauth.Claims)
	return c, ok
}

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if strings.HasPrefix(h, "Bearer ") {
				if claims, err := coreauth.ParseToken(cfg.JWTSecret, strings.TrimPrefix(h, "Bearer ")); err == nil {
					r = r.WithContext(WithClaims(r.Context(), claims))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
