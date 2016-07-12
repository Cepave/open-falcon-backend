package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/open-falcon/judge/cron"
	"github.com/open-falcon/judge/g"
	"github.com/open-falcon/judge/http"
	"github.com/open-falcon/judge/rpc"
	"github.com/open-falcon/judge/store"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(*cfg)
	logruslog.Init()

	g.InitRedisConnPool()
	g.InitHbsClient()

	store.InitHistoryBigMap()

	go http.Start()
	go rpc.Start()

	go cron.SyncStrategies()
	go cron.CleanStale()

	select {}
}
