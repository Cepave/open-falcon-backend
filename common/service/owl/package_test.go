package owl

import (
	"testing"

	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	ch "gopkg.in/check.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}
