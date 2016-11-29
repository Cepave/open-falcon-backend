package nqm

import (
	"gopkg.in/go-playground/validator.v9"
)

var Validator = validator.New()

func init() {
	Validator.RegisterValidation("nonZeroId", validNonZeroId)
}

func validNonZeroId(field validator.FieldLevel) bool {
	return field.Field().Int() != 0
}
