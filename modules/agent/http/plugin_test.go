package http

import (
	"github.com/toolkits/sys"
	"os/exec"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestRunCmdFamilyWithTimeout(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t1 := time.Now()
	cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
	err, isTimeout := cmdSessionRunWithTimeout(cmd, time.Duration(3)*time.Second)
	t2 := time.Now()
	t.Log("Time spent should less than 4 second: ", t2.Sub(t1))
	t.Log("error message is:", err)
	t.Log("Timeout happens: ", isTimeout)
	time.Sleep(3 * time.Second)
	if err != nil || isTimeout != true {
		t.Error("failed in cmdSessionRunWithTimeout")
	}
	t.Log("verify with: ps aux | head -1;ps aux| grep watch")
	t.Log("There should not be any 'watch date' existing.")

	cmd2 := exec.Command("sleep", "1")
	err2, isTimeout2 := cmdSessionRunWithTimeout(cmd2, time.Duration(3)*time.Second)
	t.Log("error message is:", err2)
	t.Log("Timeout happens: ", isTimeout2)
	if isTimeout2 != false {
		t.Error("failed in cmdSessionRunWithTimeout")
	}

	cmd3 := exec.Command("/bin/sh", "-c", "watch ls > cmd3.txt")
	cmd3.Start()
	err3, isTimeout3 := sys.CmdRunWithTimeout(cmd3, time.Duration(3)*time.Second)
	t.Log("error message is:", err3)
	t.Log("Timeout happens: ", isTimeout3)
	t.Log("verify with: ps aux | head -1;ps aux| grep watch")
	t.Log("There should be 'watch ls' existing")
}
