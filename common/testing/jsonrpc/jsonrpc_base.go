package jsonrpc

import (
	"flag"
	"fmt"
	"net/rpc"
	"testing"
	"time"

	tknet "github.com/toolkits/net"
	check "gopkg.in/check.v1"
)

var jsonRpcHost = flag.String("jsonrpc.host", "localhost", "Host of JSON-RPC")
var jsonRpcPort = flag.Int("jsonrpc.port", 80, "Port of JSON-RPC")

// Callback used to use an opened client and safty-close
type FuncJsonRpcClientCallback func(*rpc.Client)

func OpenClient(c *check.C, callback FuncJsonRpcClientCallback) {
	c.Logf("JSONRPC Connection: %s", getTargetAddress())

	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second * 3)
	c.Assert(err, check.IsNil)

	defer client.Close()

	callback(client)
}

func OpenClientBenchmark(b *testing.B, callback FuncJsonRpcClientCallback) {
	b.Logf("JSONRPC Connection: %s", getTargetAddress())

	client, err := tknet.JsonRpcClient("tcp", getTargetAddress(), time.Second * 3)
	if err != nil {
		b.Fatalf("Open TCP to address[%s] has error: %v", getTargetAddress(), err)
	}

	defer client.Close()

	callback(client)
}

var finalAddress string = ""
func getTargetAddress() string {
	if finalAddress == "" {
		finalAddress = fmt.Sprintf("%s:%d", *jsonRpcHost, *jsonRpcPort)
	}

	return finalAddress
}
