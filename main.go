package main

import (
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	InitGeneralConfig()
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
