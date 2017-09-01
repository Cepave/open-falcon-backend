package http

import (
	ch "gopkg.in/check.v1"
	"testing"

	"github.com/Cepave/open-falcon-backend/common/testing/http/gock"

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

var mockMySqlApi = gock.GockConfigBuilder.NewConfig(
	"ack.com.cc", 22060,
)
