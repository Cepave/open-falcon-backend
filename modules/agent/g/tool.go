package g

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/toolkits/file"
)

func GetCurrPluginVersion() (str string, err error) {
	if !Config().Plugin.Enabled {
		str = "0"
		err = errors.New("plugin not enabled")
		return
	}

	pluginDir := Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		str = "0"
		err = errors.New("plugin dir not existent")
		return
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		str = "0"
		return
	}

	// err is nil
	str = strings.TrimSpace(out.String())
	return
}

func GetCurrGitRepo() string {
	if !Config().Plugin.Enabled {
		return "plugin not enabled"
	}

	pluginDir := Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		return "plugin dir not existent"
	}

	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error:%s", err.Error())
	}

	return strings.TrimSpace(out.String())
}
