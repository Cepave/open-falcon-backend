package http

import (
	ch "gopkg.in/check.v1"
	"testing"

	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
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

var testFlags = tFlag.NewTestFlags()

var mockMySqlApi = gock.GockConfigBuilder.NewConfig(
	"ack.com.cc", 22060,
)

var skipItOnMySqlApi = tFlag.BuildSkipFactory(tFlag.F_ItWeb, tFlag.FeatureHelpString(tFlag.F_ItWeb))
