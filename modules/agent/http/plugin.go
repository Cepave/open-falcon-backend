package http

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/plugins"
	"github.com/toolkits/file"
)

func deleteAndCloneRepo(w http.ResponseWriter) {
	dir := g.Config().Plugin.Dir
	parentDir := file.Dir(dir)
	cmd1 := exec.Command("rm", "-rf", file.Basename(dir))
	cmd1.Dir = parentDir
	err1 := cmd1.Run()
	if err1 != nil {
		w.Write([]byte(fmt.Sprintf("\nremove the git plugin dir:%s fail. error: %s", file.Basename(dir), err1)))
		return
	} else {
		w.Write([]byte("\nremove the git plugin dir success"))
	}
	cmd2 := exec.Command("git", "clone", g.Config().Plugin.Git, file.Basename(dir))
	cmd2.Dir = parentDir
	err2 := cmd2.Run()
	if err2 != nil {
		w.Write([]byte(fmt.Sprintf("\ngit clone in dir:%s fail. error: %s", parentDir, err2)))
		return
	}
	w.Write([]byte("\ngit clone success"))
}

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
			err := cmd.Run()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("git pull in dir:%s fail. error: %s", dir, err)))
				deleteAndCloneRepo(w)
				return
			}
		} else {
			// git clone
			cmd := exec.Command("git", "clone", g.Config().Plugin.Git, file.Basename(dir))
			cmd.Dir = parentDir
			err := cmd.Run()
			if err != nil {
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
