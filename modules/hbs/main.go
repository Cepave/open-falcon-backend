package main

import (
	"flag"
	"fmt"
	"github.com/Cepave/hbs/cache"
	"github.com/Cepave/hbs/db"
	"github.com/Cepave/hbs/g"
	"github.com/Cepave/hbs/http"
	"github.com/Cepave/hbs/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	db.Init()
	cache.Init()

	go cache.DeleteStaleAgents()

	go http.Start()
	go rpc.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		db.DB.Close()
		os.Exit(0)
	}()

	select {}
}
