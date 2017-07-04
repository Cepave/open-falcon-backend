package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/fe/cache"
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/graph"
	"github.com/Cepave/open-falcon-backend/modules/fe/grpc"
	"github.com/Cepave/open-falcon-backend/modules/fe/http"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/portal"
	"github.com/Cepave/open-falcon-backend/modules/fe/model"
	"github.com/Cepave/open-falcon-backend/modules/fe/mq"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/logger"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// parse config
	vipercfg.Load()
	if err := g.ParseConfig(vipercfg.Config().GetString("config")); err != nil {
		log.Fatalln(err)
	}
	logruslog.Init()

	conf := g.Config()
	logger.SetLevelWithDefault(g.Config().Log, "info")

	model.InitDatabase()
	cache.InitCache()

	if conf.Grpc.Enabled {
		graph.Start()
		go grpc.Start()
	}
	if conf.Mq.Enabled {
		go mq.Start()
	}
	if conf.Http.Enabled {
		go http.Start()
		go portal.CornDaemonStart()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		if sig.String() == "^C" {
			os.Exit(3)
		}
	}
}
