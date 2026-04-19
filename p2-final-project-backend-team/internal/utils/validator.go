package utils

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// Register custom validation for UserRole enum
	v.RegisterValidation("userrole", validateUserRole)

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateUserRole checks if the role is a valid enum value
func validateUserRole(fl validator.FieldLevel) bool {
	role, ok := fl.Field().Interface().(enums.UserRole)
	if !ok {
		return false
	}

	// Check if role matches any valid enum value
	validRoles := []enums.UserRole{
		enums.Admin,
		enums.Student,
		enums.HeadOfDepartment,
		enums.ViceDean,
		enums.Dean,
		enums.CommitteeChair,
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}

	return false
}

// ValidateStruct is a convenience function for quick validation
func ValidateStruct(i interface{}) error {
	v := NewValidator()
	return v.Validate(i)
}
