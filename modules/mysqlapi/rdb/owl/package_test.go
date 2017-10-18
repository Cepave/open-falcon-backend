package owl

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

var inTx func(sql ...string)

var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacade()

	if DbFacade != nil {
		inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx
	}
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})

var (
	itFeatures    = tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_MySql, itSkipMessage)
)
