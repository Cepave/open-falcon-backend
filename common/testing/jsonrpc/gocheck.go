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

// Prepends "BeforeEach()" for skipping if there is no value for JSONRPC flags.
//
// Deprecated: Try to use "flag.SkipFactory"
//
// See "common/testing/flag"
func HasJsonRpcClient(c *check.C) bool {
	hasJsonRpcClient := getTestFlags().HasJsonRpcClient()

	if !hasJsonRpcClient {
		c.Skip(flagMessage)
	}

	return hasJsonRpcClient
}
