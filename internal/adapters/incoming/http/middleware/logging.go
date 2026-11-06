package middleware

import (
	"net/http"
	"time"

	"Go_Commerce/internal/infrastructure/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// LoggingMiddleware logs every request through the configured logger.
func LoggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			if rec.status == 0 {
				rec.status = http.StatusOK
			}
			args := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"bytes", rec.bytes,
				"remote_addr", r.RemoteAddr,
			}

			switch {
			case rec.status >= http.StatusInternalServerError:
				log.Error("http request", args...)
			case rec.status >= http.StatusBadRequest:
				log.Warn("http request", args...)
			default:
				log.Info("http request", args...)
			}
		})
	}
}
