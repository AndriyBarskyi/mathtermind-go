package middleware

import (
	"net/http"
	"time"

	"mathtermind-go/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// MetricsMiddleware collects request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process the request
		next.ServeHTTP(ww, r)

		// Record metrics
		duration := time.Since(start)
		logger.NewDevelopment().Info("Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", duration,
		)
	})
}

// AddMiddleware adds all our custom middleware to the router
func AddMiddleware(r *chi.Mux) {
	// Add request ID and logging
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Add metrics
	r.Use(MetricsMiddleware)

	// Add timeouts
	r.Use(middleware.Timeout(60 * time.Second))
}
