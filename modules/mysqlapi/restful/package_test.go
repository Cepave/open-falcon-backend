package restful

import (
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
var httpClientConfig = &tHttp.SlingClientConf{tHttp.NewHttpClientConfigByFlag()}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

func itSkipForGocheck(c *ch.C) {
	if !testFlags.HasHttpClient() ||
		!testFlags.HasMySql() {
		c.Skip(itSkipMessage)
	}
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
})

var (
	itFeatures    = tFlag.F_HttpClient | tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_HttpClient|tFlag.F_MySql, itSkipMessage)
	testFlags     = tFlag.NewTestFlags()
)
