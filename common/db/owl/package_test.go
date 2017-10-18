package owl

import (
	"testing"

	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
	ch "gopkg.in/check.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

var ginkgoDb = &dbTest.GinkgoDb{}

var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
})

var (
	itFeatures    = tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_MySql, itSkipMessage)
	testFlags     = tFlag.NewTestFlags()
)
