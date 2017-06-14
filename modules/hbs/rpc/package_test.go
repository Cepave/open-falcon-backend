package rpc

import (
	"flag"
	"testing"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	ch "gopkg.in/check.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

var DbFacade *f.DbFacade = db.DbFacade
var ginkgoJsonRpc = &tJsonRpc.GinkgoJsonRpc{}

func init() {
	flag.Parse()
}
