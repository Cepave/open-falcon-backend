package rpc

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"time"

	coModel "github.com/open-falcon/common/model"

	sjson "github.com/bitly/go-simplejson"
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
		response   = coModel.SimpleRpcResponse{}
		ts         *httptest.Server
		receiveCnt int64
	)

	BeforeEach(func() {
		ts = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			j, err := sjson.NewFromReader(r.Body)
			Expect(err).To(BeNil())
			arr, err := j.Array()
			Expect(err).To(BeNil())

			rowsAffectedCnt := int64(len(arr))
			defer func() {
				GinkgoT().Log("Mock rowsAffectedCnt:", rowsAffectedCnt)
				receiveCnt = rowsAffectedCnt
			}()

			resp := `{}`
			w.Write([]byte(resp))
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

		Eventually(func() int64 {
			return receiveCnt
		}, time.Second*4, time.Second/2).Should(Equal(int64(1)))
	})
}))
