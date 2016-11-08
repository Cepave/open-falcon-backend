package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/check.v1"
)

// Asserts that there is an field error with specifying name
func AssertSingleErrorForField(c *check.C, err error, expectedFieldName string) {
	c.Assert(err, check.NotNil, check.Commentf("The expected error of validation is nil"))

	validateErrors, ok := err.(validator.ValidationErrors)

	c.Assert(
		ok, check.Equals, true,
		check.Commentf("Cannot cast error to validator.ValidationErrors"),
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
