package owl

import (
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	owlDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/owl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var ginkgoDb = &tDb.GinkgoDb{}
var dbFacade = &f.DbFacade{}

var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()
	owlDb.DbFacade = dbFacade
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
	owlDb.DbFacade = dbFacade
})

func inTx(sql ...string) {
	dbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}
