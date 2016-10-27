package jsonrpc

import (
	"flag"
	"fmt"
	"net/rpc"
	tknet "github.com/toolkits/net"
	check "gopkg.in/check.v1"
	"time"
)

var jsonRpcHost = flag.String("jsonrpc.host", "localhost", "Host of JSON-RPC")
var jsonRpcPort = flag.Int("jsonrpc.port", 80, "Port of JSON-RPC")

// Callback used to use an opened client and safty-close
type FuncJsonRpcClientCallback func(*rpc.Client)

func OpenClient(c *check.C, callback FuncJsonRpcClientCallback) {
	var address = fmt.Sprintf("%s:%d", *jsonRpcHost, *jsonRpcPort)

	c.Logf("JSONRPC Connection: %s", address)

	client, err := tknet.JsonRpcClient("tcp", address, time.Second * 3)
	c.Assert(err, check.IsNil)

	defer func() {
		client.Close()
	}()

	callback(client)
}
