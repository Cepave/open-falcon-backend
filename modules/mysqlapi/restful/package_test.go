package restful

import (
	"flag"
	ch "gopkg.in/check.v1"
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dbFacade *f.DbFacade
var httpClientConfig = tHttp.NewHttpClientConfigByFlag()

func init() {
	flag.Parse()
}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
})
