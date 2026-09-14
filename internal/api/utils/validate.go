package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorMessage(err error) []string {
	fieldErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	messages := make([]string, 0, len(fieldErrors))
	for _, fieldError := range fieldErrors {
		field := fieldError.Error()
		tag := fieldError.Tag()
		param := fieldError.Param()

		var message string
		switch tag {
		case "required":
			message = field + " is required"
		case "min":
			message = field + ": length/value must be at least " + param
		case "max":
			message = field + ": length/value must be at most " + param
		case "alphanum":
			message = field + " must contain only latin letters and digits"
		case "containsany":
			message = field + " must contain at least one of: " + param
		default:
			message = field + " is invalid"
		}

		messages = append(messages, message)
	}

	return messages
}

func ValidationErrorText(err error) string {
	return strings.Join(ValidationErrorMessage(err), "; ")
}