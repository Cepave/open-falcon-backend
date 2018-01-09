package rpc

import (
	"testing"

	tflag "github.com/Cepave/open-falcon-backend/common/testing/flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var (
	jsonRpcSkipper = tflag.BuildSkipFactory(
		tflag.F_JsonRpcClient,
		tflag.FeatureHelpString(tflag.F_JsonRpcClient),
	)
)
