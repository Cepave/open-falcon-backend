package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/fe/cache"
	"github.com/freedomkk-qfeng/fe/g"
	"github.com/freedomkk-qfeng/fe/http"
	"github.com/open-falcon/fe/model"
	"github.com/toolkits/logger"
	"log"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// parse config
	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalln(err)
	}

	logger.SetLevelWithDefault(g.Config().Log, "info")

	model.InitDatabase()
	cache.InitCache()

	http.Start()
}
