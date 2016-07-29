package cron

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/plugins"
	log "github.com/Sirupsen/logrus"
)

func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {

	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		log.Println("Agent.MinePlugins content is:", resp)
		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp
		plugins.GitRepo = resp.GitRepo

		if resp.GitUpdate {
			addr := fmt.Sprintf("http://127.0.0.1%s/plugin/update", g.Config().Http.Listen)
			log.Println("GitUpdate API address is: ", addr)
			apiResp, _ := http.Get(addr)
			log.Println(&apiResp)
		}

		if g.Config().Debug {
			log.Println(&resp)
		}

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
		}

		desiredAll := make(map[string]*plugins.Plugin)

		for _, p := range pluginDirs {
			underOneDir := plugins.ListPlugins(strings.Trim(p, "/"))
			for k, v := range underOneDir {
				desiredAll[k] = v
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)

	}
}
