package rpc

import (
	"flag"
	"testing"

	tJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
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

var ginkgoJsonRpc = &tJsonRpc.GinkgoJsonRpc{}
var MOCK_URL = "localhost:5566"

func init() {
	flag.Parse()
}
