package apierror

import (
	"net/http"
)

func NewInternalError(err error) *Error {
	return &Error{
		HttpCode: http.StatusInternalServerError,
		Message:  ErrorMessageInternal.String(),
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		HttpCode: http.StatusNotFound,
		Message:  message,
	}
}

func NewBadRequestError(message string) *Error {
	return &Error{
		HttpCode: http.StatusBadRequest,
		Message:  message,
	}
}

func NewUnprocessableEntityError(message string) *Error {
	return &Error{
		HttpCode: http.StatusUnprocessableEntity,
		Message:  message,
	}
}
