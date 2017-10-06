package graph

import (
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"
	tflag "github.com/Cepave/open-falcon-backend/common/testing/flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itSkip = buildSkipper()
)

var ginkgoDb = &tDb.GinkgoDb{}
var DbFacade *f.DbFacade

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func inTx(sql ...string) {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacadeByFlag(tflag.OWL_DB_GRAPH)
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})

func buildSkipper() tflag.SkipFactory {
	httpSkipper := tflag.BuildSkipFactory(tflag.F_HttpClient, tflag.FeatureHelpString(tflag.F_HttpClient))
	dbSkipper := tflag.BuildSkipFactoryOfOwlDb(tflag.OWL_DB_GRAPH, tflag.OwlDbHelpString(tflag.OWL_DB_GRAPH))

	return httpSkipper.Compose(dbSkipper)
}
