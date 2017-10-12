package validator

import (
	ch "gopkg.in/check.v1"
	"gopkg.in/go-playground/validator.v9"
)

// Asserts that there is an field error with specifying name
func AssertSingleErrorForField(c *ch.C, err error, expectedFieldName string) {
	c.Assert(err, ch.NotNil, ch.Commentf("The expected error of validation is nil"))

	validateErrors, ok := err.(validator.ValidationErrors)

	c.Assert(
		ok, ch.Equals, true,
		ch.Commentf("Cannot cast error to validator.ValidationErrors"),
	)

	for _, fieldError := range validateErrors {
		if fieldError.Field() == expectedFieldName {
			return
		}
	}

	c.Fatalf(
		"Field[%s] should be invalidated. Errors: %v",
		expectedFieldName, validateErrors,
	)
}
