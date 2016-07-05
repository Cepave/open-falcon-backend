package receiver

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/receiver/rpc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
