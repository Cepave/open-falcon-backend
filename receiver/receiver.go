package receiver

import (
	"github.com/Cepave/transfer/receiver/rpc"
	"github.com/Cepave/transfer/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
