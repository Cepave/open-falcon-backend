package http

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/plugins"
	"github.com/Cepave/open-falcon-backend/modules/agent/session"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
	"net/http"
	"os/exec"
	"time"
)

func configPluginRoutes() {
	http.HandleFunc("/plugin/update", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Plugin.Enabled {
			w.Write([]byte("plugin not enabled"))
			return
		}

		dir := g.Config().Plugin.Dir
		parentDir := file.Dir(dir)
		file.InsureDir(parentDir)

		if file.IsExist(dir) {
			// git pull
			cmd := exec.Command("git", "pull")
			cmd.Dir = dir
			err, isTimeout := session.CmdSessionRunWithTimeout(cmd, time.Duration(600)*time.Second)
			if isTimeout {
				// has be killed
				if err == nil {
					log.Warnln("timeout and kill process git pull successfully")
				}
				if err != nil {
					log.Errorln("kill process git pull occur error:", err)
				}
				return
			}
			if err != nil {
				log.Errorf("git pull in dir: %s fail. error: %s", dir, err)
				w.Write([]byte(fmt.Sprintf("git pull in dir:%s fail. error: %s", dir, err)))
				return
			}
		} else {
			// git clone
			cmd := exec.Command("git", "clone", g.Config().Plugin.Git, file.Basename(dir))
			cmd.Dir = parentDir
			err, isTimeout := session.CmdSessionRunWithTimeout(cmd, time.Duration(600)*time.Second)
			if isTimeout {
				// has be killed
				if err == nil {
					log.Warnln("timeout and kill process git clone successfully")
				}
				if err != nil {
					log.Errorln("kill process git clone occur error:", err)
				}
				return
			}
			if err != nil {
				log.Errorf("git clone in dir:%s fail. error: %s", parentDir, err)
				w.Write([]byte(fmt.Sprintf("git clone in dir:%s fail. error: %s", parentDir, err)))
				return
			}
		}

		w.Write([]byte("success"))
	})

	http.HandleFunc("/plugin/reset", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Plugin.Enabled {
			w.Write([]byte("plugin not enabled"))
			return
		}

		dir := g.Config().Plugin.Dir

		if file.IsExist(dir) {
			cmd := exec.Command("git", "reset", "--hard")
			cmd.Dir = dir
			err := cmd.Run()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("git reset --hard in dir:%s fail. error: %s", dir, err)))
				return
			}
		}
		w.Write([]byte("success"))
	})

	http.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		//TODO: not thread safe
		RenderDataJson(w, plugins.Plugins)
	})
}
