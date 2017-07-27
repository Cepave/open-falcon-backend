package jsonrpc

import (
	"testing"
	"time"

	tknet "github.com/toolkits/net"
	check "gopkg.in/check.v1"
)

func OpenClient(c *check.C, callback FuncJsonRpcClientCallback) {
	c.Logf("JSONRPC Connection: %s", getTargetAddress())

	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second*3)
	c.Assert(err, check.IsNil)

	defer client.Close()

	callback(client)
}

func OpenClientBenchmark(b *testing.B, callback FuncJsonRpcClientCallback) {
	b.Logf("JSONRPC Connection: %s", getTargetAddress())

	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second*3)
	if err != nil {
		b.Fatalf("Open TCP to address[%s] has error: %v", getTargetAddress(), err)
	}

	defer client.Close()

	callback(client)
}

func HasJsonRpcServ(c *check.C) bool {
	var hasJsonRpcHost = *jsonRpcHost != ""

	if !hasJsonRpcHost {
		c.Skip("Skip json-rpc testing. Needs \"-jsonrpc.host=<Host address of JSON-RPC>\"")
	}

	return hasJsonRpcHost
}
