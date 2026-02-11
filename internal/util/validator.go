package util

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func FormatValidationErrors(err error) ValidationErrorResponse {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: getErrorMessage(e),
			})
		}
	}

	return ValidationErrorResponse{Errors: errors}
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short, minimum is " + e.Param()
	case "max":
		return "Value is too long, maximum is " + e.Param()
	case "e164":
		return "Invalid phone number format (use E.164 format, e.g., +14155552671)"
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	default:
		return "Invalid value for " + e.Tag()
	}
}
