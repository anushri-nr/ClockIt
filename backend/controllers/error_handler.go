package controllers

import (
	"github.com/go-playground/validator/v10"
)

// GetValidationErrors extracts detailed validation error messages from Gin binding errors
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			var message string
			switch tag {
			case "required":
				message = fieldName + " is required"
			case "email":
				message = fieldName + " must be a valid email address"
			case "min":
				message = fieldName + " must be at least " + fieldError.Param() + " characters"
			case "max":
				message = fieldName + " must not exceed " + fieldError.Param() + " characters"
			case "pattern":
				message = fieldName + " format is invalid"
			default:
				message = fieldName + " validation failed: " + tag
			}

			errors[fieldName] = message
		}
	}

	if len(errors) == 0 {
		errors["error"] = "request body is invalid or malformed"
	}

	return errors
}
