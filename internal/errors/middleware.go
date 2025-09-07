package errors

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"mathtermind-go/internal/logger"
)

// HandlerFunc is a custom http.HandlerFunc that returns an error
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Middleware converts our HandlerFunc to a standard http.Handler with error handling
func Middleware(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add request ID to context if not present
		if GetReqID(r.Context()) == "" {
			ctx := context.WithValue(r.Context(), requestIDKey, generateRequestID())
			r = r.WithContext(ctx)
		}

		// Set request ID in response header
		reqID := GetReqID(r.Context())
		if reqID != "" {
			w.Header().Set("X-Request-ID", reqID)
		}

		// Call the handler function
		if err := h(w, r); err != nil {
			// Log the error with request ID
			logger.NewDevelopment().Error("Request error", 
				"error", err,
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			// Write the error response to the client
			WriteError(w, err)
		}
	})
}

// Recoverer is a middleware that recovers from panics and converts them to 500 errors
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := GetReqID(r.Context())
				
				// Log the panic with request ID
				logger.NewDevelopment().Error("Recovered from panic",
					"panic", rec,
					"request_id", reqID,
				)
				
				// Create an internal server error
				err := Errorf(ErrCodeInternal, "Internal server error")
				if e, ok := rec.(error); ok {
					err = err.WithError(e)
				}
				
				// Write the error response
				WriteError(w, err)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// requestIDKey is the context key for request IDs
type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// GetReqID retrieves the request ID from context
func GetReqID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
