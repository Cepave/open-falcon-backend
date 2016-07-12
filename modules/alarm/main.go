package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cepave/alarm/cron"
	"github.com/Cepave/alarm/g"
	"github.com/Cepave/alarm/http"
	"github.com/Cepave/alarm/logger"
	"github.com/Cepave/alarm/model"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitRedisConnPool()
	model.InitDatabase()
	logger.InitLogger()

	go http.Start()
	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineSms()
	go cron.CombineMail()
	go cron.CombineQQ()
	go cron.CombineServerchan()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
