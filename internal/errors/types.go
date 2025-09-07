package errors

import "fmt"

// ErrorCode represents a unique code for each error type
type ErrorCode string

const (
	// Generic errors
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidation  ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrCodeForbidden   ErrorCode = "FORBIDDEN"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// Database errors
	ErrCodeDBError          ErrorCode = "DB_ERROR"
	ErrCodeDBConnection    ErrorCode = "DB_CONNECTION_ERROR"
	ErrCodeDBQuery         ErrorCode = "DB_QUERY_ERROR"
	ErrCodeDBMigration     ErrorCode = "DB_MIGRATION_ERROR"

	// Authentication & Authorization
	ErrCodeAuthError        ErrorCode = "AUTH_ERROR"
	ErrCodeLoginError       ErrorCode = "AUTH_LOGIN_ERROR"
	ErrCodePermissionDenied ErrorCode = "AUTH_PERMISSION_DENIED"
	ErrCodeTokenError       ErrorCode = "AUTH_TOKEN_ERROR"

	// Business logic
	ErrCodeBusinessLogic  ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrCodeInvalidState   ErrorCode = "INVALID_STATE"
	ErrCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"

	// External services
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// Error represents a custom error type with error code and details
type Error struct {
	// Machine-readable error code
	Code ErrorCode `json:"code"`
	// Human-readable message
	Message string `json:"message"`
	// Additional error details
	Details map[string]any `json:"details,omitempty"`
	// The underlying error that was returned by a function
	Err error `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// WithDetails adds key-value pairs to the error details
func (e *Error) WithDetails(details map[string]any) *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// WithError wraps an underlying error
func (e *Error) WithError(err error) *Error {
	e.Err = err
	return e
}

// New creates a new error with the given code and message
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: make(map[string]any),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]any),
	}
}

// Errorf creates a new formatted error
func Errorf(code ErrorCode, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Details: make(map[string]any),
	}
}

// Is checks if the target error is of type *Error with the same error code
func Is(err error, targetCode ErrorCode) bool {
	errObj, ok := err.(*Error)
	if !ok {
		return false
	}
	return errObj.Code == targetCode
}

// As converts an error to *Error if possible
func As(err error) (*Error, bool) {
	e, ok := err.(*Error)
	return e, ok
}
