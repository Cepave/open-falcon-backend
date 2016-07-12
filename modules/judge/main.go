package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/judge/cron"
	"github.com/Cepave/open-falcon-backend/modules/judge/g"
	"github.com/Cepave/open-falcon-backend/modules/judge/http"
	"github.com/Cepave/open-falcon-backend/modules/judge/rpc"
	"github.com/Cepave/open-falcon-backend/modules/judge/store"
	"os"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
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
