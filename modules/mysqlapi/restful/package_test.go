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

var testFlags = tFlag.NewTestFlags()
func itSkipForGocheck(c *ch.C) {
	if !testFlags.HasHttpClient() ||
		!testFlags.HasMySqlOfOwlDb(tFlag.OWL_DB_PORTAL) {
		c.Skip(itSkipMessage)
	}
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	portalDbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_PORTAL)
	bossDbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_BOSS)
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(portalDbFacade)
	ginkgoDb.ReleaseDbFacade(bossDbFacade)

	portalDbFacade = nil
	bossDbFacade = nil
})

var (
	httpMessage  = tFlag.FeatureHelpString(tFlag.F_HttpClient)
	portalMessage = tFlag.OwlDbHelpString(tFlag.OWL_DB_PORTAL)

	itSkipMessage = fmt.Sprintf(
		"%s; %s", httpMessage, portalMessage,
	)

	httpClientSkip = tFlag.BuildSkipFactory(tFlag.F_HttpClient, httpMessage)

	itSkipOnPortal = httpClientSkip.Compose(
		tFlag.BuildSkipFactoryOfOwlDb(tFlag.OWL_DB_PORTAL, portalMessage),
	)

	cmdbFlags = (tFlag.OWL_DB_PORTAL | tFlag.OWL_DB_BOSS)
	itSkipForCmdb = httpClientSkip.Compose(
		tFlag.BuildSkipFactoryOfOwlDb(
			cmdbFlags,
			tFlag.OwlDbHelpString(cmdbFlags),
		),
	)
)

var portalDbFacade *f.DbFacade
func inPortalTx(sql ...string) {
	portalDbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var bossDbFacade *f.DbFacade
func inBossTx(sql ...string) {
	bossDbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}
