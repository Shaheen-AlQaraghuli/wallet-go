package json

import (
	"github.com/gin-gonic/gin"
	"wallet/internal/util/http/apierror"
)

func SendApiValidationError(c *gin.Context, err error) {
	validationError := apierror.NewValidatorError(c, err)
	c.JSON(validationError.HttpCode, validationError)
}

func SendGenericAPIError(c *gin.Context, err error) {
	genericError := apierror.NewUnprocessableEntityError(err.Error())
	c.JSON(genericError.HttpCode, genericError)
}

func SendBadRequestError(c *gin.Context, message string) {
	badRequestError := apierror.NewBadRequestError(message)
	c.JSON(badRequestError.HttpCode, badRequestError)
}
