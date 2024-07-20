package utils

import (
	"fmt"
	"go_blogs/models"

	"github.com/go-playground/validator/v10"
)

func getValidationErrorMessageFromTag(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "is invalid email"
	case "mongodb":
		return "is invalid ID"
	case "min":
		return fmt.Sprintf("must be longer than %s", err.Param())
	case "max":
		return fmt.Sprintf("must be shorter than %s", err.Param())
	}
	return "is invalid"
}

func TransformValidationErrorFormat(errs error) []models.ValidationErrorResponse {
	var customErrors []models.ValidationErrorResponse

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			customError := models.ValidationErrorResponse{
				Field:   err.Field(),
				Value:   err.Value(),
				Tag:     err.Tag(),
				Message: fmt.Sprintf("%s %s", err.Field(), getValidationErrorMessageFromTag(err)),
			}
			customErrors = append(customErrors, customError)
		}
	}

	return customErrors
}
