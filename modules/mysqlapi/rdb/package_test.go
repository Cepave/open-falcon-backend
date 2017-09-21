package rdb

import (
	"testing"

	ch "gopkg.in/check.v1"

	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itFeatures    = tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_MySql, itSkipMessage)
	testFlags     = tFlag.NewTestFlags()
)

func itSkipForGocheck(c *ch.C) {
	if !testFlags.HasMySql() {
		c.Skip(itSkipMessage)
	}
}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

func inTx(sql ...string) {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})
