package session

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
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
func CmdSessionRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	log.Debugln("show cmd sysProcAttr.Setsid", cmd.SysProcAttr.Setsid)
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
