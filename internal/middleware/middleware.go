package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// requestIDKey is the context key for request IDs
type requestIDKey struct{}

var requestIDKeyInstance = requestIDKey{}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate request ID
		reqID := generateRequestID()

		// Add request ID to context
		ctx := context.WithValue(r.Context(), requestIDKeyInstance, reqID)
		r = r.WithContext(ctx)

		// Set request ID in response header
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r)
	})
}

// GetReqID retrieves the request ID from context
func GetReqID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKeyInstance).(string); ok {
		return reqID
	}
	return ""
}

// StructuredErrorMiddleware handles errors in a structured way
func StructuredErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				logger := slog.Default()
				reqID := GetReqID(r.Context())
				logger.Error("Panic occurred",
					"error", err,
					"request_id", reqID,
				)

				// Write error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware collects request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger := slog.Default()
		logger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// AddMiddleware adds all our custom middleware to the router
func AddMiddleware(r *chi.Mux) {
	r.Use(RequestIDMiddleware)
	r.Use(StructuredErrorMiddleware)
	r.Use(MetricsMiddleware)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
}
