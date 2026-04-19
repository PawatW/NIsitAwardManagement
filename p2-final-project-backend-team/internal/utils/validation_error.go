package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// AddUserFormatValidationError formats validation errors for AddUser request
func AddUserFormatValidationError(err error) []ValidationError {
	var errors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var message string

			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", e.Field())
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", e.Field())
			case "userrole":
				message = fmt.Sprintf("%s must be a valid role (ADMIN, STUDENT, HEAD_OF_DEPARTMENT, VICE_DEAN, DEAN, COMMITTEE_CHAIR)", e.Field())
			default:
				message = fmt.Sprintf("%s is invalid", e.Field())
			}

			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	return errors
}

// Example: Add more specific validation formatters for other operations
// func UpdateUserFormatValidationError(err error) []ValidationError { ... }
// func LoginFormatValidationError(err error) []ValidationError { ... }
