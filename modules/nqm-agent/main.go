package main

import (
	"fmt"
	"os"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"

	flag "github.com/spf13/pflag"
)

func main() {
	cfgFilePtr := flag.String("c", "cfg.json", "nqm's configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	InitGeneralConfig(*cfgFilePtr)
	logruslog.Init()
	InitRPC()

	go Query()
	Measure()

	select {}
}
