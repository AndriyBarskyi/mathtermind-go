package errors

import (
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	return "validation failed"
}

// ToMap converts validation errors to a map
func (ve ValidationErrors) ToMap() map[string]string {
	errs := make(map[string]string)
	for _, e := range ve {
		errs[e.Field] = e.Message
	}
	return errs
}

// ValidateStruct validates a struct and returns validation errors if any
func ValidateStruct(s any) error {
	validate := validator.New()
	
	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return newValidationError(ve)
		}
		
		// If it's not a validation error, return as internal error
		return Wrap(err, ErrCodeInternal, "Failed to validate request")
	}
	
	return nil
}

func newValidationError(ve validator.ValidationErrors) error {
	var validationErrs ValidationErrors
	
	for _, e := range ve {
		validationErrs = append(validationErrs, ValidationError{
			Field:   e.Field(),
			Message: getValidationMessage(e),
		})
	}
	
	err := New(ErrCodeValidation, "Validation failed")
	err.Details["errors"] = validationErrs.ToMap()
	return err
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return e.Error()
	}
}
