package apierror

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidatorError returns the parsed and formatted validation error
func NewValidatorError(c context.Context, err error) *Error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return &Error{
			HttpCode:     http.StatusBadRequest,
			Message:      err.Error(),
		}
	}

	apiErr := &Error{
		HttpCode:    http.StatusUnprocessableEntity,
		Message:      ErrorMessageValidation.String(),
		Errors: make([]ValidationError, 0),
	}

	for _, err := range errs {
		fieldErr, ok := err.(validator.FieldError)
		if !ok {
			return &Error{
				HttpCode:     http.StatusUnprocessableEntity,
				Message:      err.Error(),
			}
		}

		validationErr := ValidationError{
			Message:    fieldErr.Error(),
			Source:     fmt.Sprintf("/%s", strings.ToLower(err.Field())),
			fieldError: fieldErr,
		}

		apiErr.Errors = append(apiErr.Errors, validationErr)
	}

	return apiErr
}
