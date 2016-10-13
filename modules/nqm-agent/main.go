package main

import (
	"fmt"
	"os"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	InitJSONConfig()
	InitGeneralConfig()
	logruslog.Init()
	InitRPC()

	go Query()
	Measure()

	select {}
}
