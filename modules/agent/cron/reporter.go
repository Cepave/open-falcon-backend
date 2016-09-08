package cron

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	log "github.com/Sirupsen/logrus"
)

func ReportAgentStatus() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go reportAgentStatus(time.Duration(g.Config().Heartbeat.Interval) * time.Second)
	}
}

func reportAgentStatus(interval time.Duration) {
	for {
		hostname, err := g.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("error:%s", err.Error())
		}

		currPluginVersion, currPluginErr := g.GetCurrPluginVersion()
		if currPluginErr != nil {
			log.Warnln("GetCurrPluginVersion returns: ", currPluginErr)
		}

		currPluginRepo, currRepoErr := g.GetCurrGitRepo()
		if currRepoErr != nil {
			log.Warnln("GetCurrGitRepo returns: ", currRepoErr)
		}

		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: currPluginVersion,
			GitRepo:       currPluginRepo,
		}

		log.Debugln("show req of Agent.ReportStatus: ", req)
		var resp model.SimpleRpcResponse
		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			log.Errorln("call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
		}

		time.Sleep(interval)
	}
}
