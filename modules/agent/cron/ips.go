package cron

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	log "github.com/sirupsen/logrus"
	"time"
)

func SyncTrustableIps() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", model.NullRpcRequest{}, &ips)
		if err != nil {
			log.Errorln("call Agent.TrustableIps fail", err)
			continue
		}

		g.SetTrustableIps(ips)
	}
}
