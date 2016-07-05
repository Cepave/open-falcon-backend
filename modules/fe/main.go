package main

import (
	"flag"
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/fe/cache"
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/graph"
	"github.com/Cepave/open-falcon-backend/modules/fe/grpc"
	"github.com/Cepave/open-falcon-backend/modules/fe/http"
	"github.com/Cepave/open-falcon-backend/modules/fe/model"
	"github.com/Cepave/open-falcon-backend/modules/fe/mq"
	"github.com/toolkits/logger"
	"log"
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
	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalln(err)
	}

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
