package restful

import (
	"fmt"
	ch "gopkg.in/check.v1"
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dbFacade *f.DbFacade
var clientConfig = tHttp.NewHttpClientConfigByFlag()
var httpClientConfig = &tHttp.SlingClientConf{clientConfig}
var gentlemanClientConfig = &tHttp.GentlemanClientConf{clientConfig}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

func itSkipForGocheck(c *ch.C) {
	if !testFlags.HasHttpClient() ||
		!testFlags.HasMySqlOfOwlDb(tFlag.OWL_DB_PORTAL) {
		c.Skip(itSkipMessage)
	}
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_PORTAL)
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
})

var (
	httpMessage  = tFlag.FeatureHelpString(tFlag.F_HttpClient)
	mysqlMessage = tFlag.OwlDbHelpString(tFlag.OWL_DB_PORTAL)

	itSkipMessage = fmt.Sprintf(
		"%s; %s", httpMessage, mysqlMessage,
	)
	itSkip = tFlag.BuildSkipFactory(tFlag.F_HttpClient, httpMessage).Compose(
		tFlag.BuildSkipFactoryOfOwlDb(tFlag.OWL_DB_PORTAL, mysqlMessage),
	)
	testFlags = tFlag.NewTestFlags()
)

func inTx(sql ...string) {
	dbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}
