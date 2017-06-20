package cron

import (
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/plugins"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
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

func dirFilter(userBindDir []string) (resultDir []string) {
	// remove reserved dir: data
	for _, val := range userBindDir {
		if !strings.HasPrefix(val, "data/") && val != "data" {
			resultDir = append(resultDir, val)
		}
	}
	return
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
			log.Errorln("call Agent.MinePlugins fail", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}
		log.Debugln("Response of RPC call Agent.MinePlugins: ", &resp)

		pluginDirs = dirFilter(resp.Plugins)
		timestamp = resp.Timestamp

		// git repo updating.
		log.Debugln("GitRepo auto update with HBS: ", g.Config().Plugin.AutoGitRepoUpdate)
		if g.Config().Plugin.AutoGitRepoUpdate {
			if currPluginRepo, currRepoErr := plugins.GetCurrGitRepo(); currRepoErr != nil {
				log.Warnln("GetCurrGitRepo returns: ", currRepoErr)
				if !file.IsExist(g.Config().Plugin.Dir) {
					log.Debugln("local git repo not existent.")
					log.Debugln("initializing git repo by HBS.")
					plugins.UpdatePlugin("", resp.GitRepo)
				}
			} else {
				if currPluginRepo != resp.GitRepo {
					log.Debugln("local git repo != HBS's git repo.")
					log.Debugln("git remote set-url origin", resp.GitRepo)
					plugins.SetCurrGitRepo(resp.GitRepo)
					plugins.UpdatePlugin("", resp.GitRepo)
				}
			}
		}

		// git commit sync
		log.Debugln("Git commits auto sync: ", g.Config().Plugin.AutoGitUpdate)
		if g.Config().Plugin.AutoGitUpdate {
			if currPluginRepo, currRepoErr := plugins.GetCurrGitRepo(); currRepoErr != nil {
				log.Warnln("GetCurrGitRepo returns: ", currRepoErr)
				if !file.IsExist(g.Config().Plugin.Dir) {
					log.Debugln("local git repo not existent.")
					log.Debugln("initializing git repo by config.")
					plugins.UpdatePlugin("", "")
				}
			} else {
				if hash, err := plugins.GitLsRemote(currPluginRepo, "refs/heads/master"); err != nil {
					log.Warnln("Error retrieving git-repo:", currPluginRepo, err)
				} else {
					log.Debugln("Get newest plugin hash from: ", currPluginRepo, hash)
					if currHash, currErr := plugins.GetCurrPluginVersion(); currErr != nil {
						log.Warnln("GetCurrPluignVersion returns: ", currHash)
					} else {
						if currHash != hash {
							log.Debugln("local git's HEAD != origin's HEAD.")
							log.Debugln("git fetch; git reset --hard ", hash)
							plugins.UpdatePlugin(hash, currPluginRepo)
						}
					}
				}
			}
		}

		if g.Config().Debug {
			log.Println(&resp)
			log.Println(pluginDirs)
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
