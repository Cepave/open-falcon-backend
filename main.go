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

	for {
		go func() {
			probingCmd, targets, agentPtr, err := QueryTask()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Execution the probing command:", probingCmd[0:9])

			rawData := Probe(probingCmd)
			jsonParams := MarshalIntoParameters(rawData, targets, agentPtr)
			Push(jsonParams)
		}()
		time.Sleep(time.Minute)
	}
	select {}
}
