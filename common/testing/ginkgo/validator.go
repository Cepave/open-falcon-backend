package ginkgo

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"gopkg.in/go-playground/validator.v9"
)

func MatchFieldErrorOnName(fieldName string) types.GomegaMatcher {
	return &fieldErrorOnNameMatcher{
		fieldName: fieldName,
	}
}

type fieldErrorOnNameMatcher struct {
	fieldName string
}

func (matcher *fieldErrorOnNameMatcher) Match(actual interface{}) (success bool, err error) {
	validateErrors, ok := actual.(validator.ValidationErrors)

	if !ok {
		return false, fmt.Errorf("Type of error(%t) is not type of \"validator.ValidationErrors\"", actual)
	}

	if validateErrors == nil {
		return false, fmt.Errorf("Value of error is nil")
	}

	for _, fieldError := range validateErrors {
		if fieldError.Field() == matcher.fieldName {
			return true, nil
		}
	}

	return false, nil
}

func (matcher *fieldErrorOnNameMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Expected field name for validation error\n\t[ %s ]\nCurrent errors:\n\t%v",
		matcher.fieldName, actual,
	)
}

func (matcher *fieldErrorOnNameMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Not expected field name for validation error\n\t[ %s ]\nCurrent errors:\n\t%v",
		matcher.fieldName, actual,
	)
}
