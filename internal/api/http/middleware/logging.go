package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs the request details.
func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			log.Info("HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.String("status_text", http.StatusText(ww.Status())),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
