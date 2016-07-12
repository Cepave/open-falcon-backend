package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/agent/cron"
	"github.com/Cepave/open-falcon-backend/modules/agent/funcs"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"github.com/Cepave/open-falcon-backend/modules/agent/http"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(*cfg)
	logruslog.Init()

	g.InitRootDir()
	g.InitPublicIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()
	cron.Collect()

	go http.Start()

	select {}

}
