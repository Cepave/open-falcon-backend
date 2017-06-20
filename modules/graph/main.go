package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cepave/open-falcon-backend/modules/graph/api"
	"github.com/Cepave/open-falcon-backend/modules/graph/g"
	"github.com/Cepave/open-falcon-backend/modules/graph/http"
	"github.com/Cepave/open-falcon-backend/modules/graph/index"
	"github.com/Cepave/open-falcon-backend/modules/graph/rrdtool"
)

func start_signal(pid int, cfg *g.GlobalConfig) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			if cfg.Http.Enabled {
				http.Close_chan <- 1
				<-http.Close_done_chan
			}
			log.Println("http stop ok")

			if cfg.Rpc.Enabled {
				api.Close_chan <- 1
				<-api.Close_done_chan
			}
			log.Println("rpc stop ok")

			rrdtool.Out_done_chan <- 1
			rrdtool.FlushAll(true)
			log.Println("rrdtool stop ok")

			log.Println(pid, "exit")
			os.Exit(0)
		}
	}
}

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if vipercfg.Config().GetBool("vg") {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()
	// init db
	g.InitDB()
	// rrdtool before api for disable loopback connection
	rrdtool.Start()
	// start api
	go api.Start()
	// start indexing
	index.Start()
	// start http server
	go http.Start()

	start_signal(os.Getpid(), g.Config())
}
