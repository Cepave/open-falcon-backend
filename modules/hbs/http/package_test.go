package http

import (
	"testing"

	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itEnabled = tHttp.GinkgoHttpIt.NeedItWeb
	itConfig  = tHttp.NewHttpClientConfigByFlag()
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Suite")
}
