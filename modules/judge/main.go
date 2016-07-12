package main

import (
	"fmt"
	"os"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/judge/cron"
	"github.com/Cepave/open-falcon-backend/modules/judge/g"
	"github.com/Cepave/open-falcon-backend/modules/judge/http"
	"github.com/Cepave/open-falcon-backend/modules/judge/rpc"
	"github.com/Cepave/open-falcon-backend/modules/judge/store"
	flag "github.com/spf13/pflag"
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
