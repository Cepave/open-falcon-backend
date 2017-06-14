package rpc

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"

	apiModel "github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	coModel "github.com/open-falcon/common/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test rpc call: Agent.ReportStatus", ginkgoJsonRpc.NeedJsonRpc(func() {
	var (
		request = coModel.AgentReportRequest{
			Hostname:      "test-g-01",
			IP:            "123.45.61.81",
			AgentVersion:  "4.5.31",
			PluginVersion: "1.2.12",
		}
		response = coModel.SimpleRpcResponse{}
		ts       *httptest.Server
	)

	BeforeEach(func() {
		ts = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			decorder := json.NewDecoder(r.Body)
			var rAgents []*apiModel.FalconAgentHeartbeat
			err := decorder.Decode(&rAgents)
			Expect(err).To(BeNil())
			defer r.Body.Close()

			rowsAffectedCnt := int64(len(rAgents))
			res := apiModel.FalconAgentHeartbeatResult{rowsAffectedCnt}
			resp, err := json.Marshal(res)
			Expect(err).To(BeNil())
			w.Write(resp)
		}))
		l, err := net.Listen("tcp", MOCK_URL)
		Expect(err).To(BeNil())
		ts.Listener = l
		ts.Start()

		GinkgoT().Logf("Mock server at: %s", ts.URL)
		GinkgoT().Log("Please set 'mysql_api.host' in hbs to addr of mock server")
	})

	AfterEach(func() {
		ts.Close()
	})

	It("should get correct value", func() {
		ginkgoJsonRpc.OpenClient(func(client *rpc.Client) {
			err := client.Call("Agent.ReportStatus", request, &response)

			Expect(err).To(BeNil())
		})
	})
}))
