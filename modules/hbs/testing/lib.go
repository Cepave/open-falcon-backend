package testing

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	tknet "github.com/toolkits/net"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type RpcTestEnv struct {
	Port      int
	RpcClient *rpc.Client

	stop bool
	wait chan bool
}

func DefaultListenAndExecute(
	rcvr interface{}, rpcCallback func(*RpcTestEnv),
) {
	var rpcTestEnvInstance = RpcTestEnv{
		Port: 18080,
		stop: false,
		wait: make(chan bool, 1),
	}

	rpcTestEnvInstance.ListenAndExecute(
		rcvr, rpcCallback,
	)
}
func (rpcTestEnvInstance *RpcTestEnv) ListenAndExecute(
	rcvr interface{}, rpcCallback func(*RpcTestEnv),
) {
	server := rpc.NewServer()
	server.Register(rcvr)

	var address = fmt.Sprintf("localhost:%d", rpcTestEnvInstance.Port)

	/**
	 * Listening RPC
	 */
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("Cannot listen TCP[%d]. Error: %v", rpcTestEnvInstance.Port, err))
	}

	log.Printf("Listen RPC at port [%d]", rpcTestEnvInstance.Port)
	// :~)

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				if rpcTestEnvInstance.stop {
					log.Printf("Stop RPC server")
					rpcTestEnvInstance.wait <- true
				}

				break
			}

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()

	/**
	 * Initialize RPC client
	 */
	rpcTestEnvInstance.RpcClient, err = tknet.JsonRpcClient("tcp", address, time.Second*3)
	if err != nil {
		panic(fmt.Errorf("Initialize RPC client error: %v", err))
	}
	log.Printf("Construct RPC Client")
	// :~)

	defer func() {
		rpcTestEnvInstance.RpcClient.Close()
		rpcTestEnvInstance.stop = true
		listener.Close()
		<-rpcTestEnvInstance.wait
	}()

	rpcCallback(rpcTestEnvInstance)
}
