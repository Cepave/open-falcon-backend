package validate

import (
	"testing"

	v "gopkg.in/go-playground/validator.v9"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func newValidator() *v.Validate {
	validator := v.New()
	RegisterDefaultValidators(validator)
	return validator
}
