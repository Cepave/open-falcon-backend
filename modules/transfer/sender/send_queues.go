package sender

import (
	"github.com/Cepave/transfer/g"
	nlist "github.com/toolkits/container/list"
)

func initSendQueues() {
	cfg := g.Config()
	for node, _ := range cfg.Judge.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	if cfg.Tsdb.Enabled {
		TsdbQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}

	Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	InfluxdbQueues["default"] = Q

	if cfg.NqmRpc.Enabled {
		NqmRpcQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}
