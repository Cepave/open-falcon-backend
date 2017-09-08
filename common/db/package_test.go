package db

import (
	"testing"

	ch "gopkg.in/check.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}
