package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/sender/cron"
	"github.com/Cepave/open-falcon-backend/modules/sender/g"
	"github.com/Cepave/open-falcon-backend/modules/sender/http"
	"github.com/Cepave/open-falcon-backend/modules/sender/redis"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"
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
