package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	onet "github.com/Cepave/open-falcon-backend/common/net"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
)

func StartRpc() {
	if !g.Config().Rpc.Enabled {
		logger.Info("RPC(JSON) service is disabled")
		return
	}

	address := g.Config().Rpc.Listen
	logger.Infof("Initializes RPC(JSON) service: %s", address)

	listener := onet.MustInitTcpListener(address)
	listenerCtrl := onet.NewListenerController(listener)
	defer listenerCtrl.Close()

	server := rpc.NewServer()
	server.Register(new(Transfer))

	listenerCtrl.AcceptLoop(
		func(conn net.Conn) {
			server.ServeCodec(jsonrpc.NewServerCodec(conn))
		},
	)
}
