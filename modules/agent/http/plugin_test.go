package http

import (
	"os/exec"
	"testing"
	"time"

	"github.com/toolkits/sys"
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
	err := cmd.Start()
	if err != nil {
		t.Log("cmd start fails: ", err)
	}
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(3)*time.Second)
	t2 := time.Now()
	t.Log("Time spent should less than 4 second: ", t2.Sub(t1))
	t.Log("error message is:", err)
	t.Log("Timeout happens: ", isTimeout)
	time.Sleep(3 * time.Second)
	if err != nil || isTimeout != true {
		t.Error("failed in sys.CmdRunWithTimeout")
	}

	cmd2 := exec.Command("sleep", "1")
	err2 := cmd2.Start()
	if err2 != nil {
		t.Log("cmd start fails: ", err2)
	}
	err2, isTimeout2 := sys.CmdRunWithTimeout(cmd, time.Duration(3)*time.Second)
	t.Log("error message is:", err2)
	t.Log("Timeout happens: ", isTimeout2)
	if isTimeout2 != false {
		t.Error("failed in sys.CmdRunWithTimeout")
	}
}
