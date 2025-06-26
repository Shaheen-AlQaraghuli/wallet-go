package apierror

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMessage string

const (
	ErrorMessageValidation ErrorMessage = "validation error"
	ErrorMessageInternal ErrorMessage = "something went wrong"
)

func (code ErrorMessage) String() string {
	return string(code)
}

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Source     string `json:"source"`
	Message    string `json:"message"`
	fieldError validator.FieldError
}

// Error represents a standardized error structure
type Error struct {
	HttpCode         int                `json:"-"`
	Message          string             `json:"message,omitempty"`
	Errors []ValidationError  `json:"errors,omitempty"`
}

func (e Error) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("http_code= %d, message: %s", e.HttpCode, e.Message))

	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf(", field: %s, error: %s", err.Source, err.Message))
	}

	return sb.String()
}


