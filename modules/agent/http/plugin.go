package http

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/plugins"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func sessionKill(cmd *exec.Cmd) error {
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		log.Println("failed to kill processes of the same session: ", err)
		return err
	}
	return nil
}

// cmdSessionRunWithTimeout runs cmd, and kills the children on timeout
func cmdSessionRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		log.Errorln(cmd.Path, " start fails: ", err)
		return err, false
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(timeout):
		log.Printf("timeout, process:%s will be killed", cmd.Path)
		err := sessionKill(cmd)
		return err, true
	case err := <-done:
		return err, false
	}
}

func DeleteAndCloneRepo(pluginDir string, gitRemoteAddr string) (out string) {
	parentDir := file.Dir(pluginDir)

	absPath, _ := filepath.Abs(pluginDir)
	if absPath == "/" {
		out = fmt.Sprintf("\nRemove directory:%s is not allowed.", absPath)
		log.Errorln(out)
		return
	}
	err1 := os.RemoveAll(file.Basename(pluginDir))
	if err1 != nil {
		out = fmt.Sprintf("\nremove the git plugin dir:%s fail. error: %s", file.Basename(pluginDir), err1)
		log.Errorln(out)
		return
	} else {
		out = fmt.Sprintf("\nremove the git plugin dir:%s success", file.Basename(pluginDir))
	}

	cmd := exec.Command("git", "clone", gitRemoteAddr, file.Basename(pluginDir))
	cmd.Dir = parentDir
	err2, isTimeout := cmdSessionRunWithTimeout(cmd, time.Duration(600)*time.Second)
	if isTimeout {
		// has be killed
		if err2 == nil {
			log.Warnln("timeout and kill process git clone successfully")
		}
		if err2 != nil {
			log.Errorln("kill process git clone occur error:", err2)
		}
		return
	}
	if err2 != nil {
		log.Errorln("git clone run fails:", err2)
		out = out + fmt.Sprintf("\ngit clone in dir:%s fail. error: %s", parentDir, err2)
		return
	}
	out = out + "\ngit clone success"
	return
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
			err, isTimeout := cmdSessionRunWithTimeout(cmd, time.Duration(600)*time.Second)
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
				w.Write([]byte(DeleteAndCloneRepo(dir, plugins.GitRepo)))
				return
			}
		} else {
			// git clone
			cmd := exec.Command("git", "clone", plugins.GitRepo, file.Basename(dir))
			cmd.Dir = parentDir
			err, isTimeout := cmdSessionRunWithTimeout(cmd, time.Duration(600)*time.Second)
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
