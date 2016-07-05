package main

import (
	"flag"
	"fmt"
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
	versionGit := flag.Bool("vg", false, "show version and git commit log")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// config
	g.ParseConfig(*cfg)
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
		conf.ReadConf("./conf/lambdaSetup.json")
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
