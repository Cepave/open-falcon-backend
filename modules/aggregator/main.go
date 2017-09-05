package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cepave/open-falcon-backend/common/sdk/sender"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/cron"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/db"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/g"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	db.Init()

	go http.Start()
	go cron.UpdateItems()

	// sdk configuration
	sender.Debug = g.Config().Debug
	sender.PostPushUrl = g.Config().Api.PushApi

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
