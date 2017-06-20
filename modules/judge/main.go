package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/judge/cron"
	"github.com/Cepave/open-falcon-backend/modules/judge/g"
	"github.com/Cepave/open-falcon-backend/modules/judge/http"
	"github.com/Cepave/open-falcon-backend/modules/judge/rpc"
	"github.com/Cepave/open-falcon-backend/modules/judge/store"
	log "github.com/sirupsen/logrus"
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
	g.InitLastEvents()

	store.InitHistoryBigMap()

	supervisorChn := make(chan string)

	go http.Start(supervisorChn)
	go rpc.Start(supervisorChn)

	go cron.SyncStrategies(supervisorChn)
	go cron.CleanStale(supervisorChn)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		if sig.String() == "^C" {
			os.Exit(3)
		}
	//keep all routine can auto recovery from any painc actions
	case sup := <-supervisorChn:
		if sup == "http" {
			log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
			go http.Start(supervisorChn)
		} else if sup == "rpc" {
			log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
			go rpc.Start(supervisorChn)
		} else if sup == "SyncStrategies" {
			log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
			go cron.SyncStrategies(supervisorChn)
		} else if sup == "CleanStale" {
			log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
			go cron.SyncStrategies(supervisorChn)
		} else {
			log.Fatalf("got worng params of supervisorChn -> %v .", sup)
		}
	}
}
