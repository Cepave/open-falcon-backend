package gin

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/conform"
	"gopkg.in/go-playground/validator.v9"
)

// Used to be handled globally
type ValidationError struct {
	errors validator.ValidationErrors
}

// Implements error interface
func (err ValidationError) Error() string {
	return err.errors.Error()
}

// Conforms the object and then validates the object,
// if any error occurs, panic with validation error.
//
// See leebeson/conform: https://github.com/leebenson/conform
// See go-playground/validator: https://godoc.org/gopkg.in/go-playground/validator.v9
func ConformAndValidateStruct(object interface{}, v *validator.Validate) {
	conform.MustConform(object)

	err := v.Struct(object)

	if err == nil {
		return
	}

	validatorErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		panic(fmt.Errorf("Unknown validation error: %v", err))
	}

	panic(ValidationError{validatorErrors})
}
