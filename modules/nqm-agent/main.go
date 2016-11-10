package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/fsnotify/fsnotify"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	InitConfig()
	logruslog.Init()

	hbsTicker = time.NewTicker(Config().Hbs.Interval * time.Second)
	hbsTickerUpdated = make(chan bool)

	vipercfg.Config().WatchConfig()
	vipercfg.Config().OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		InitConfig()
		logruslog.Init()
		hbsTickerUpdated <- true
		GenMeta()
		InitRPC()
	})

	GenMeta()
	InitRPC()

	go Query()
	Measure()

	select {}
}
