package nqm

import (
	"testing"

	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ch "gopkg.in/check.v1"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

type dbTestSuite struct{}

func (s *dbTestSuite) SetUpSuite(c *ch.C) {
	qtest.InitDb(c)
	initServices()
}
func (s *dbTestSuite) TearDownSuite(c *ch.C) {
	qtest.ReleaseDb(c)
}

var ginTestServer = tHttp.GinTestServer
