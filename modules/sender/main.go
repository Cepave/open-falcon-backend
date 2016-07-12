package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/sender/cron"
	"github.com/Cepave/open-falcon-backend/modules/sender/g"
	"github.com/Cepave/open-falcon-backend/modules/sender/http"
	"github.com/Cepave/open-falcon-backend/modules/sender/redis"
	flag "github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"
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

	vipercfg.Load()
	g.ParseConfig(*cfg)
	logruslog.Init()
	cron.InitWorker()
	redis.InitConnPool()

	go http.Start()
	go cron.ConsumeSms()
	go cron.ConsumeMail()
	go cron.ConsumeQQ()
	go cron.ConsumeServerchan()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		redis.ConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
