package validate

import (
	v "gopkg.in/go-playground/validator.v9"
)

func RegisterDefaultValidators(validator *v.Validate) {
	validator.RegisterValidation(TagNonZeroSlice, NonZeroSlice)
}
