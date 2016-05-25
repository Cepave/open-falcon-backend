package main

import (
	"log"
	"time"
)

func tcpconn() {
	for {
		func() {
			probingCmd, targets, agentPtr, err := makeTasks("tcpconn")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Execution of the probing command:", probingCmd)

			rawData := Probe(probingCmd, "tcpconn")
			/*for _, row := range rawData {
				fmt.Println(row)
			}*/

			jsonParams := MarshalIntoParameters(rawData, targets, agentPtr, "tcpconn")
			/*
				for i, _ := range jsonParams {
					println(jsonParams[i].String())
					println("===")
				}
			*/
			Push(jsonParams)
		}()

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
			probingCmd, targets, agentPtr, err := makeTasks("fping")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Execution of the probing command:", probingCmd[0:9])

			rawData := Probe(probingCmd, "fping")
			jsonParams := MarshalIntoParameters(rawData, targets, agentPtr, "fping")
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
