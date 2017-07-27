package restful

import (
	"net/http"

	json "github.com/Cepave/open-falcon-backend/common/json"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test /health", ginkgoDb.NeedDb(func() {
	It("returns the JSON data", func() {
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewSlingByBase().
				Get("health"),
		)
		jsonBody := resp.GetBodyAsJson()
		GinkgoT().Logf("[Mysql API Module Response] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
	})
}))
