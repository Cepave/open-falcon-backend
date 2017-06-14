package jsonrpc

import (
	"time"

	tknet "github.com/toolkits/net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type GinkgoJsonRpc struct{}

func (g *GinkgoJsonRpc) OpenClient(callback FuncJsonRpcClientCallback) {
	GinkgoT().Logf("JSONRPC Connection: %s", getTargetAddress())

	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second*3)
	Expect(err).To(BeNil())

	defer client.Close()

	callback(client)
}

func (g *GinkgoJsonRpc) NeedJsonRpc(src func()) func() {
	return func() {
		BeforeEach(func() {
			if *jsonRpcFlag == false {
				Skip("Skip json-rpc testing. Needs \"-jsonrpc\"")
			}
		})

		src()
	}
}
