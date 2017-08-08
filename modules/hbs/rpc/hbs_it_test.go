package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test rpc call [Hbs.GetExpressions]", ginkgoJsonRpc.NeedJsonRpc(func() {
	request := commonModel.NullRpcRequest{}
	response := commonModel.ExpressionResponse{}
	var ts *httptest.Server
	tsResp := `
  [
    {
      "id":3,
      "metric":"ss.close.wait",
      "tags":{
        "endpoint":"oth-bj-119-090-062-121"
      },
      "func":"all(#1)",
      "operator":"!=",
      "right_value":0,
      "max_step":1,
      "priority":4,
      "note":"boss oth-bj-119-090-062-121 连接数大于10",
      "action_id":91
    }
  ]
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
			err := client.Call("Hbs.GetExpressions", request, &response)
			respStr := fmt.Sprintf("%v", response)
			GinkgoT().Logf("RPC Response(%.5000s)", respStr)
			Expect(err).To(BeNil())
			Expect(response.Expressions).To(HaveLen(1))
		})
	})
}))

var _ = Describe("Test rpc call [Hbs.GetStrategies]", ginkgoJsonRpc.NeedJsonRpc(func() {
	request := commonModel.NullRpcRequest{}
	response := commonModel.StrategiesResponse{}
	var ts *httptest.Server
	tsResp := `
  [
     {
        "hostname":"",
        "strategies":[
           {
              "id":5,
              "metric":"net.if.in.bits",
              "tags":"iface=eth0",
              "func":"all(#1)",
              "operator":"!=",
              "right_value":"0",
              "maxStep":1,
              "priority":4,
              "note":"boss oth-bj-119-090-062-121 连接数大于10",
              "tpl":{
                 "id":6,
                 "name":"traffic warning",
                 "parent_id":2,
                 "action_id":3,
                 "creator":"test-user"
              }
           }
        ]
     }
  ]
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
			err := client.Call("Hbs.GetStrategies", request, &response)
			respStr := fmt.Sprintf("%v", response)
			GinkgoT().Logf("RPC Response(%.5000s)", respStr)
			Expect(err).To(BeNil())
			Expect(response.HostStrategies).To(HaveLen(1))
		})
	})
}))
