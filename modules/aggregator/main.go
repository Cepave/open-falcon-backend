package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cepave/open-falcon-backend/modules/aggregator/cron"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/db"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/g"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/http"
	"github.com/Cepave/open-falcon-backend/sdk/graph"
	"github.com/Cepave/open-falcon-backend/sdk/portal"
	"github.com/Cepave/open-falcon-backend/sdk/sender"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if vipercfg.Config().GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()
	db.Init()

	go http.Start()
	go cron.UpdateItems()

	// sdk configuration
	graph.GraphLastUrl = g.Config().Api.GraphLast
	sender.Debug = g.Config().Debug
	sender.PostPushUrl = g.Config().Api.Push
	portal.HostnamesUrl = g.Config().Api.Hostnames

	sender.StartSender()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
