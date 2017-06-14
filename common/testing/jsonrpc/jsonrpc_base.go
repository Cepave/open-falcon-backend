package jsonrpc

import (
	"flag"
	"fmt"
	"net/rpc"
)

var jsonRpcFlag = flag.Bool("jsonrpc", false, "Enable JSON-RPC")
var jsonRpcHost = flag.String("jsonrpc.host", "localhost", "Host of JSON-RPC")
var jsonRpcPort = flag.Int("jsonrpc.port", 80, "Port of JSON-RPC")

// Callback used to use an opened client and safty-close
type FuncJsonRpcClientCallback func(*rpc.Client)

var finalAddress string = ""

func getTargetAddress() string {
	if finalAddress == "" {
		finalAddress = fmt.Sprintf("%s:%d", *jsonRpcHost, *jsonRpcPort)
	}

	return finalAddress
}
