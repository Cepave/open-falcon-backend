package jsonrpc

import (
	"time"

	tknet "github.com/toolkits/net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type GinkgoJsonRpc struct{}

func (g *GinkgoJsonRpc) OpenClient(callback FuncJsonRpcClientCallback) {
	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second*3)
	Expect(err).To(Succeed(), "Cannot open json rpc client target: %s. Error: %v", getTargetAddress(), err)
	defer client.Close()

	callback(client)
}

// Prepends "BeforeEach()" for skipping if there is no value for JSONRPC flags.
//
// Deprecated: Try to use "flag.SkipFactory"
//
// See "common/testing/flag"
func (g *GinkgoJsonRpc) NeedJsonRpc(src func()) func() {
	return func() {
		BeforeEach(func() {
			if !getTestFlags().HasJsonRpcClient() {
				Skip(flagMessage)
			}
		})

		src()
	}
}
