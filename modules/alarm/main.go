package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/alarm/cron"
	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	"github.com/Cepave/open-falcon-backend/modules/alarm/http"
	"github.com/Cepave/open-falcon-backend/modules/alarm/model"
	"github.com/spf13/pflag"
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
	g.InitRedisConnPool()
	model.InitDatabase()

	go http.Start()
	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineSms()
	go cron.CombineMail()
	go cron.CombineQQ()
	go cron.CombineServerchan()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
