package main

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"os"

	"github.com/Cepave/open-falcon-backend/modules/task/collector"
	"github.com/Cepave/open-falcon-backend/modules/task/g"
	"github.com/Cepave/open-falcon-backend/modules/task/http"
	"github.com/Cepave/open-falcon-backend/modules/task/index"
	"github.com/Cepave/open-falcon-backend/modules/task/proc"
)

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
	// proc
	proc.Start()

	// graph index
	index.Start()
	// collector
	collector.Start()

	// http
	http.Start()

	select {}
}
