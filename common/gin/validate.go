package gin

import (
	"fmt"
	"github.com/leebenson/conform"
	"gopkg.in/go-playground/validator.v9"
)

// Used to be handled globally
type ValidationError struct {
	errors validator.ValidationErrors
}

func (err ValidationError) Error() string {
	return err.errors.Error()
}

// Conforms the object and then validates the object,
// if any error occurs, panic with validation error.
func ConformAndValidateStruct(object interface{}, v *validator.Validate) {
	conform.Strings(object)

	err := v.Struct(object)

	if err == nil {
		return
	}

	validatorErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		panic(fmt.Errorf("Unknown validation error: %v", err))
	}

	panic(ValidationError{ validatorErrors })
}
