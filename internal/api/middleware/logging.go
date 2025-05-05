package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs the request details.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			log.Printf(
				"[%s] %s %s %d %s %dB in %s",
				r.Method,
				r.RemoteAddr,
				r.URL.Path,
				ww.Status(),
				http.StatusText(ww.Status()),
				ww.BytesWritten(),
				time.Since(start),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
