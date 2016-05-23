package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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

	go QueryHbs()

	go func() {
		for {
			func() {
				probingCmd, targets, agentPtr, err := makeTasks()
				if err != nil {
					log.Println(err)
					return
				}
				probingCmd = []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a", "127.0.0.1"}
				log.Println("Execution the probing command:", probingCmd[0:9])

				rawData := Probe(probingCmd)
				jsonParams := MarshalIntoParameters(rawData, targets, agentPtr)
				Push(jsonParams)
			}()

			dur := time.Second * time.Duration(GetGeneralConfig().Agent.FpingInterval)
			time.Sleep(dur)
		}
	}()
	select {}
}
