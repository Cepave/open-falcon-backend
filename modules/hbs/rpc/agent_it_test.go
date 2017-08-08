package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"time"

	coModel "github.com/Cepave/open-falcon-backend/common/model"

	sjson "github.com/bitly/go-simplejson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

var _ = Describe("Test rpc call [Agent.BuiltinMetrics]", ginkgoJsonRpc.NeedJsonRpc(func() {
	request := coModel.AgentHeartbeatRequest{
		Hostname: "cnc-he-060-008-151-208",
		Checksum: "d94ad826797905118ca30ba11a0273ad4303a8ca",
	}
	response := coModel.BuiltinMetricResponse{}
	var ts *httptest.Server
	tsResp := `
	{
	   "metrics":[
	      {
	         "metric":"net.if.in.bits",
	         "tags":"iface=eth0"
	      }
	   ],
	   "checksum":"e2c569be17396eca2a2e3c11578123ed",
	   "timestamp":1501491450
	}
  `

	BeforeEach(func() {
		ts = httptest.NewUnstartedServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tsResp))
			}))
		l, err := net.Listen("tcp", MOCK_URL)
		Expect(err).To(BeNil())
		ts.Listener = l
		ts.Start()

		GinkgoT().Logf("Mock server at: %s", ts.URL)
	})

	AfterEach(func() {
		ts.Close()
	})

	It("should get correct value", func() {
		ginkgoJsonRpc.OpenClient(func(client *rpc.Client) {
			err := client.Call("Agent.BuiltinMetrics", request, &response)
			respStr := fmt.Sprintf("%v", response)
			GinkgoT().Logf("RPC Response(%.5000s)", respStr)
			Expect(err).To(BeNil())
			Expect(response.Checksum).To(Equal("e2c569be17396eca2a2e3c11578123ed"))
			Expect(response.Timestamp).To(Equal(int64(1501491450)))
			Expect(response.Metrics).To(HaveLen(1))
		})
	})
}))

var _ = Describe("[Intg] Test rpc call: Agent.MinePlugins", ginkgoJsonRpc.NeedJsonRpc(func() {

	DescribeTable("when parameter is",
		func(request coModel.AgentHeartbeatRequest, expectedPluginNum int) {
			response := &coModel.AgentPluginsResponse{}
			ginkgoJsonRpc.OpenClient(func(client *rpc.Client) {
				err := client.Call("Agent.MinePlugins", request, &response)
				GinkgoT().Logf("RPC Response(%v)", response)
				Expect(err).To(BeNil())
				Expect(response.Plugins).To(HaveLen(expectedPluginNum))
			})
		},
		Entry("Nil, should get nil plugins", nil, 0),
		Entry("Not found but valide value, shold get empty plugins",
			coModel.AgentHeartbeatRequest{
				"test-agent-mineplugins",
				"Not-checked-value",
			}, 0,
		),
		Entry("Found and valid value, should get its plugins",
			coModel.AgentHeartbeatRequest{
				"cnc-he-060-008-151-208",
				"Not-checked-value",
			}, 8,
		),
	)
}))
