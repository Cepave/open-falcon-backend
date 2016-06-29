package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	cfgFilePtr := flag.String("c", "cfg.json", "nqm's configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	InitGeneralConfig(*cfgFilePtr)
	InitRPC()

	go Query()
	Measure()

	select {}
}
