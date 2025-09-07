package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler is an HTTP middleware that handles errors
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				handlePanic(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ErrorResponse represents the JSON structure for error responses
type ErrorResponse struct {
	Error ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// WriteError writes an error response in JSON format
func WriteError(w http.ResponseWriter, err error) {
	errResp := toErrorResponse(err)
	statusCode := statusCodeFromError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: errResp}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func toErrorResponse(err error) ErrorPayload {
	if customErr, ok := err.(*Error); ok {
		return ErrorPayload{
			Code:    string(customErr.Code),
			Message: customErr.Message,
			Details: customErr.Details,
		}
	}

	// Handle standard errors
	return ErrorPayload{
		Code:    string(ErrCodeInternal),
		Message: "Internal server error",
		Details: map[string]interface{}{
			"original_error": err.Error(),
		},
	}
}

func statusCodeFromError(err error) int {
	switch e := err.(type) {
	case *Error:
		switch e.Code {
		case ErrCodeValidation:
			return http.StatusBadRequest
		case ErrCodeUnauthorized:
			return http.StatusUnauthorized
		case ErrCodeForbidden:
			return http.StatusForbidden
		case ErrCodeNotFound:
			return http.StatusNotFound
		case ErrCodeDBError, ErrCodeDBConnection, ErrCodeDBQuery, ErrCodeDBMigration:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}

func handlePanic(w http.ResponseWriter, r interface{}) {
	err := Errorf(ErrCodeInternal, "Internal server error")
	err.Details["panic"] = r
	WriteError(w, err)
}
