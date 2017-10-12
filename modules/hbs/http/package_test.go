package http

import (
	"testing"

	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itConfig         = tHttp.NewHttpClientConfigByFlag()
	fakeServerConfig = &tHttp.FakeServerConfig{"127.0.0.1", 6040}
	features         = tFlag.F_HttpClient | tFlag.F_ItWeb
	sf               = tFlag.BuildSkipFactory(features, tFlag.FeatureHelpString(features))
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Suite")
}

var testFlags = tFlag.NewTestFlags()
