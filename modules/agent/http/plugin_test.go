package http

import (
	"os/exec"
	"testing"
	"time"
)

func TestDeleteAndCloneRepo(t *testing.T) {
	out := DeleteAndCloneRepo("./plugin", "https://github.com/humorless/openfalcon-plugin.git")
	t.Log("test deleteAndCloneRepo: ", out)
	out = DeleteAndCloneRepo("/", "https://github.com/humorless/openfalcon-plugin.git")
	t.Log("test deleteAndCloneRepo: ", out)
}

func TestRunCmdWithTimeout(t *testing.T) {
	cmd := exec.Command("sleep", "500")
	t1 := time.Now()
	RunCmdWithTimeout(cmd, 3)
	t2 := time.Now()
	t.Log("Time spent should less than 4 second: ", t2.Sub(t1))
}
