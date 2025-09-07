package errors

// Common error constructors

// NotFound creates a new not found error
func NotFound(resource string, id interface{}) *Error {
	err := New(ErrCodeNotFound, "Resource not found")
	err.Details["resource"] = resource
	err.Details["id"] = id
	return err
}

// Unauthorized creates a new unauthorized error
func Unauthorized(message string) *Error {
	if message == "" {
		message = "You are not authorized to perform this action"
	}
	return New(ErrCodeUnauthorized, message)
}

// Forbidden creates a new forbidden error
func Forbidden(message string) *Error {
	if message == "" {
		message = "You don't have permission to access this resource"
	}
	return New(ErrCodeForbidden, message)
}

// Validation creates a new validation error
func Validation(message string, details map[string]interface{}) *Error {
	err := New(ErrCodeValidation, message)
	if details != nil {
		err.Details = details
	}
	return err
}

// Internal creates a new internal server error
func Internal(message string, err error) *Error {
	if message == "" {
		message = "An internal error occurred"
	}
	return New(ErrCodeInternal, message).WithError(err)
}

// Database creates a new database error
func Database(err error) *Error {
	return New(ErrCodeDBError, "Database error").WithError(err)
}

// BadRequest creates a new bad request error
func BadRequest(message string) *Error {
	if message == "" {
		message = "Invalid request"
	}
	return New(ErrCodeValidation, message)
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code ErrorCode, message string) *Error {
	return Wrap(err, code, message)
}

// IsError checks if an error is of a specific error code
func IsError(err error, code ErrorCode) bool {
	return Is(err, code)
}

// GetErrorCode returns the error code from an error
func GetErrorCode(err error) ErrorCode {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return ""
}
