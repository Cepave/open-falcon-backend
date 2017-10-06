package graph

import (
	"testing"

	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var ginkgoDb = &tDb.GinkgoDb{}

func inTx(sql ...string) {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_GRAPH)
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})

var (
	itSkipMessage = tFlag.OwlDbHelpString(tFlag.OWL_DB_GRAPH)
	itSkip        = tFlag.BuildSkipFactoryOfOwlDb(tFlag.OWL_DB_GRAPH, itSkipMessage)
)
