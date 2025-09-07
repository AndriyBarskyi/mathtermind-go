package errors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mathtermind-go/internal/errors"
)

func TestErrorHandling(t *testing.T) {
	handler := errors.Middleware(func(w http.ResponseWriter, r *http.Request) error {
		// Example of validation error
		if r.URL.Query().Get("id") == "" {
			return errors.Validation("Invalid input", map[string]any{
				"id": "ID is required",
			})
		}

		// Example of not found error
		if r.URL.Query().Get("id") == "999" {
			return errors.NotFound("user", "999")
		}

		// Example of internal error
		if r.URL.Query().Get("id") == "error" {
			return errors.Internal("Failed to process request", errors.New(errors.ErrCodeDBError, "database connection failed"))
		}

		// Success case
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
		return nil
	})

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "success",
			url:            "/test?id=123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "validation error",
			url:            "/test",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			url:            "/test?id=999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			url:            "/test?id=error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
