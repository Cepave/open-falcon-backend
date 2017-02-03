package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/http"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	oos "github.com/Cepave/open-falcon-backend/common/os"
	"os"
	"syscall"
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

	db.Init()
	cache.Init()

	go cache.DeleteStaleAgents()

	go http.Start()
	go rpc.Start()

	oos.HoldingAndWaitSignal(
		func(signal os.Signal) {
			rpc.Stop()
			db.Release()
		},
		os.Interrupt, os.Kill,
		syscall.SIGTERM,
	)
}
