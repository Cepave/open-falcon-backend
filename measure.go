package main

import (
	"log"
	"time"
)

func tcpconn() {
	for {
		log.Println("Execution of tcpcpnn")

		dur := time.Second * time.Duration(GetGeneralConfig().Agent.TcpconnInterval)
		time.Sleep(dur)
	}
}

func tcpping() {
	for {
		log.Println("Execution of tcpping")

		dur := time.Second * time.Duration(GetGeneralConfig().Agent.TcppingInterval)
		time.Sleep(dur)
	}
}

func fping() {
	for {
		func() {
			probingCmd, targets, agentPtr, err := makeTasks()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Execution of the probing command:", probingCmd[0:9])

			rawData := Probe(probingCmd)
			jsonParams := MarshalIntoParameters(rawData, targets, agentPtr)
			Push(jsonParams)
		}()

		dur := time.Second * time.Duration(GetGeneralConfig().Agent.FpingInterval)
		time.Sleep(dur)
	}
}

func Measure() {
	go fping()
	go tcpping()
	go tcpconn()
}
