package owl

import (
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var httpClient = tHttp.GentlemanClientConf{tHttp.NewHttpClientConfigByFlag()}

var ginkgoDb = &tDb.GinkgoDb{}

var inTx func(sql ...string)

var dbFacade *f.DbFacade

var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()

	if dbFacade != nil {
		inTx = dbFacade.SqlDbCtrl.ExecQueriesInTx
	}
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
})
