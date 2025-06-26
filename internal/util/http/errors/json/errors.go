package json

import (
	"wallet/internal/util/http/apierror"
	"github.com/gin-gonic/gin"
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
