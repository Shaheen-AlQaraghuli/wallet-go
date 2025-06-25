package types

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

func RegisterValidations() error {
	if err := registerEnumValidation("transactionStatusEnum", GetTransactionStatuses()); err != nil {
		return err
	}

	if err := registerEnumValidation("transactionTypeEnum", GetTransactionTypes()); err != nil {
		return err
	}

	if err := registerEnumValidation("currencyEnum", GetCurrencies()); err != nil {
		return err
	}

	if err := registerEnumValidation("walletStatusEnum", GetWalletStatuses()); err != nil {
		return err
	}

	return nil
}

func registerEnumValidation[T comparable](tag string, allValues []T) error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("validator engine is not of type *validator.Validate")
	}

	return validate.RegisterValidation(
		tag,
		func(fl validator.FieldLevel) bool {
			val, ok := fl.Field().Interface().(T)

			if !ok {
				return false
			}

			return slices.Contains(allValues, val)
		},
	)
}
