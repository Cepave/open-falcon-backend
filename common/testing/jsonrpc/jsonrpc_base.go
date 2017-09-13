package jsonrpc

import (
	"fmt"
	"net/rpc"

	tflag "github.com/Cepave/open-falcon-backend/common/testing/flag"
)

// Callback used to use an opened client and safety-close
type FuncJsonRpcClientCallback func(*rpc.Client)

var testFlags *tflag.TestFlags

func getTestFlags() *tflag.TestFlags {
	if testFlags == nil {
		testFlags = tflag.NewTestFlags()
	}

	return testFlags
}

var finalAddress string = ""

func getTargetAddress() string {
	host, port := getTestFlags().GetJsonRpcClient()

	if finalAddress == "" {
		finalAddress = fmt.Sprintf("%s:%d", host, port)
	}

	return finalAddress
}

var flagMessage = fmt.Sprintf("Skip MySql Test: -owl.test=%s", tflag.FeatureHelp(tflag.F_MySql)[0])
