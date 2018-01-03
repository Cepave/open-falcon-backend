package receiver

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/receiver/rpc"
)

func Start() {
	go rpc.StartRpc()
}
