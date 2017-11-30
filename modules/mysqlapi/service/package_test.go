package service

import (
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var (
	ginkgoDb = &tDb.GinkgoDb{}
	dbFacade = &f.DbFacade{}

	itFeatures    = tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_MySql, itSkipMessage)
)

var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()
	rdb.DbFacade = dbFacade
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
	rdb.DbFacade = dbFacade
})
