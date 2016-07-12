package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/fe/cache"
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/graph"
	"github.com/Cepave/open-falcon-backend/modules/fe/grpc"
	"github.com/Cepave/open-falcon-backend/modules/fe/http"
	"github.com/Cepave/open-falcon-backend/modules/fe/model"
	"github.com/Cepave/open-falcon-backend/modules/fe/mq"
	log "github.com/Sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/toolkits/logger"
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

	// parse config
	vipercfg.Load()
	if err := g.ParseConfig(*cfg); err != nil {
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
	}

	select {}
}
