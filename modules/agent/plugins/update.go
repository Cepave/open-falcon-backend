package plugins

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

const zeroHash string = "0"

func GetCurrPluginVersion() (str string, err error) {
	if !g.Config().Plugin.Enabled {
		str = zeroHash
		err = errors.New("plugin not enabled")
		return
	}

	pluginDir := g.Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		str = zeroHash
		err = errors.New("plugin dir not existent")
		return
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		str = zeroHash
		return
	}

	// err is nil
	str = strings.TrimSpace(out.String())
	return
}

func GetCurrGitRepo() (str string, err error) {
	if !g.Config().Plugin.Enabled {
		str = zeroHash
		err = errors.New("plugin not enabled")
		return
	}

	pluginDir := g.Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		str = zeroHash
		err = errors.New("plugin dir not existent")
		return
	}

	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		str = zeroHash
		return
	}

	// err is nil
	str = strings.TrimSpace(out.String())
	return
}

// borrow the source code from satori
var updateInflight bool = false
var lastPluginUpdate int64 = 0

func UpdatePlugin(ver string, repo string) error {
	cfg := g.Config().Plugin

	if !cfg.Enabled {
		log.Debugln("plugin not enabled, not updating")
		return fmt.Errorf("plugin not enabled")
	}

	if updateInflight {
		log.Debugln("Previous update inflight, do nothing")
		return nil
	}

	// TODO: add to config
	if time.Now().Unix()-lastPluginUpdate < 600 {
		log.Debugln("Previous update too recent, do nothing")
		return nil
	}

	parentDir := file.Dir(cfg.Dir)
	file.InsureDir(parentDir)

	if ver == "" {
		ver = "origin/master"
	}
	if repo == "" {
		repo = cfg.Git
	}

	var buf bytes.Buffer

	if file.IsExist(cfg.Dir) {
		// git fetch
		log.Println("Begin update plugins by fetch")
		updateInflight = true
		defer func() { updateInflight = false }()
		lastPluginUpdate = time.Now().Unix()

		buf.Reset()
		cmd := exec.Command("timeout", "60s", "git", "fetch")
		cmd.Dir = cfg.Dir
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err := cmd.Run()
		if err != nil {
			s := fmt.Sprintf("Update plugins by fetch error: %s. message: %s", err, buf.String())
			log.Println(s)
			return fmt.Errorf("git fetch in dir:%s fail. error: %s", cfg.Dir, err)
		}

		buf.Reset()
		cmd = exec.Command("git", "reset", "--hard", ver)
		cmd.Dir = cfg.Dir
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err = cmd.Run()
		if err != nil {
			s := fmt.Sprintf("git reset --hard failed: %s. message: %s", err, buf.String())
			log.Println(s)
			return fmt.Errorf("git reset --hard in dir:%s fail. error: %s", cfg.Dir, err)
		}
		log.Println("Update plugins by fetch complete")
	} else {
		// git clone
		log.Println("Begin update plugins by clone")
		lastPluginUpdate = time.Now().Unix()
		buf.Reset()
		cmd := exec.Command("timeout", "200s", "git", "clone", repo, file.Basename(cfg.Dir))
		cmd.Dir = parentDir
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err := cmd.Run()
		if err != nil {
			s := fmt.Sprintf("Update plugins by clone error: %s. message: %s", err, buf.String())
			log.Println(s)
			return fmt.Errorf("git clone in dir:%s fail. error: %s", parentDir, err)
		}

		buf.Reset()
		cmd = exec.Command("git", "reset", "--hard", ver)
		cmd.Dir = cfg.Dir
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err = cmd.Run()
		if err != nil {
			s := fmt.Sprintf("git reset --hard failed: %s. message: %s", err, buf.String())
			log.Println(s)
			return fmt.Errorf("git reset --hard in dir:%s fail. error: %s", cfg.Dir, err)
		}
		log.Println("Update plugins by clone complete")
	}
	return nil
}

func GitLsRemote(gitRepo string, refs string) (string, error) {
	// This function depends on git command
	if resultStr, err := sys.CmdOut("timeout", "30s", "git", "ls-remote", gitRepo, refs); err != nil {
		return "", err
	} else {
		// resultStr should be:
		// cb7a2998571cb25693867afcb24a7331f597768e        refs/heads/master
		strList := strings.Fields(resultStr)
		return strList[0], nil
	}
}

func SetCurrGitRepo(gitRepo string) error {
	cmd := exec.Command("git", "remote", "set-url", "origin", gitRepo)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = g.Config().Plugin.Dir
	e := cmd.Run()
	if e != nil {
		log.Warnln("Git remote set-url error:", e, out.String())
	} else {
		log.Debugln("Git remote set-url successfully.")
	}
	return e
}
