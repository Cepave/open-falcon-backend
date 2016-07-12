package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	flag "github.com/spf13/pflag"
	"os"
	"os/signal"

	"github.com/Cepave/open-falcon-backend/modules/query/conf"
	"github.com/Cepave/open-falcon-backend/modules/query/database"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	ginHttp "github.com/Cepave/open-falcon-backend/modules/query/gin_http"
	"github.com/Cepave/open-falcon-backend/modules/query/graph"
	"github.com/Cepave/open-falcon-backend/modules/query/grpc"
	"github.com/Cepave/open-falcon-backend/modules/query/http"
	"github.com/Cepave/open-falcon-backend/modules/query/proc"
)

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// config
	vipercfg.Load()
	g.ParseConfig(*cfg)
	logruslog.Init()
	gconf := g.Config()
	// proc
	proc.Start()

	// graph
	graph.Start()

	if gconf.Grpc.Enabled {
		// grpc
		go grpc.Start()
	}

	if gconf.GinHttp.Enabled {
		//lambdaSetup
		database.Init()
		conf.ReadConf()
		go ginHttp.StartWeb()
	}

	if gconf.Http.Enabled {
		// http
		go http.Start()
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
