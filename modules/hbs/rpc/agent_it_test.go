package rpc

import (
	"net/rpc"

	testJsonRpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"
	"github.com/open-falcon/common/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test rpc call: Agent.ReportStatus", func() {
	var (
		request = model.AgentReportRequest{
			Hostname:      "test-g-01",
			IP:            "123.45.61.81",
			AgentVersion:  "4.5.31",
			PluginVersion: "1.2.12",
		}
		response      = model.SimpleRpcResponse{}
		ginkgoJsonRpc = &testJsonRpc.GinkgoJsonRpc{}
	)

	It("should get correct value", func() {
		ginkgoJsonRpc.OpenClient(func(client *rpc.Client) {
			err := client.Call("Agent.ReportStatus", request, &response)

			Expect(err).To(BeNil())
		})
	})
})
