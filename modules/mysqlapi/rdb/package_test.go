package rdb

import (
	"flag"
	"testing"

	ch "gopkg.in/check.v1"

	tDb "github.com/Cepave/open-falcon-backend/common/testing/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
	DbFacade = ginkgoDb.InitDbFacade()
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})
