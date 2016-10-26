package plugins

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	close(this.Quit)
}

func noOwnerExecPerm(fpath string) bool {
	info, err := os.Stat(fpath)
	if err != nil {
		log.Errorln("cannot stat file", err)
	}

	perm := info.Mode().Perm()
	if (perm & 0100) != 0100 {
		return true
	} else {
		return false
	}
}

func hasShebang(fpath string) bool {
	// examine if it has shebang in file's first two characters.
	file, err := os.Open(fpath)
	if err != nil {
		log.Println("cannot open plugin script file", err)
	}
	defer file.Close()

	var line []string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = strings.Split(scanner.Text(), " ")
	}
	if strings.HasPrefix(line[0], "#!") {
		return true
	} else {
		return false
	}
}

func getInterpreterCmd(fpath string) []string {
	file, err := os.Open(fpath)
	if err != nil {
		log.Println("cannot open plugin script file", err)
	}
	defer file.Close()

	var itpr []string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		itpr = strings.Split(strings.Trim(scanner.Text(), " "), " ")
		itpr[0] = strings.TrimPrefix(itpr[0], "#!")
	}
	itpr = append(itpr, fpath)

	return itpr
}

func PluginRun(plugin *Plugin) {

	timeout := plugin.Cycle*1000 - 500
	fpath := filepath.Join(g.Config().Plugin.Dir, plugin.FilePath)

	if !file.IsExist(fpath) {
		log.Println("no such plugin:", fpath)
		return
	}

	debug := g.Config().Debug
	if debug {
		log.Println(fpath, "running...")
	}

	var cmd *exec.Cmd
	if noOwnerExecPerm(fpath) && hasShebang(fpath) {
		itprcmd := getInterpreterCmd(fpath)
		cmd = exec.Command(itprcmd[0], itprcmd[1:]...)
		if debug {
			log.Println("[INFO]", fpath, "has shebang but no owner exec perm.")
		}
	} else {
		cmd = exec.Command(fpath)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Errorln(fpath, " start fails: ", err)
	}

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	if errStr != "" {
		logFile := filepath.Join(g.Config().Plugin.LogDir, plugin.FilePath+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Errorf("write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			log.Println("[INFO] timeout and kill process", fpath, "successfully")
		}

		if err != nil {
			log.Errorln("kill process", fpath, "occur error:", err)
		}

		return
	}

	if err != nil {
		log.Errorln("exec plugin", fpath, "fail. error:", err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		if debug {
			log.Println("[DEBUG] stdout of", fpath, "is blank")
		}
		return
	}

	var metrics []*model.MetricValue
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Errorf("json.Unmarshal stdout of %s fail. error:%s stdout: \n%s\n", fpath, err, stdout.String())
		return
	}

	// fill in fields
	sec := plugin.Cycle
	now := time.Now().Unix()
	hostname, err := g.Hostname()
	if err != nil {
		hostname = ""
	}

	for j := 0; j < len(metrics); j++ {
		metrics[j].Step = int64(sec)
		metrics[j].Endpoint = hostname
		metrics[j].Timestamp = now
	}

	toTransfer, toMQ := g.DemultiplexMetrics(metrics)
	g.SendToTransfer(toTransfer)
	g.SendToMQ(toMQ)
}
