package main

import (
	"fmt"
	"log"
)

type ParamToAgent struct {
	Metric      string      `json:"metric"`
	Endpoint    string      `json:"endpoint"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	Tags        string      `json:"tags"`
	Timestamp   int64       `json:"timestamp"`
	Step        int64       `json:"step"`
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	InitGeneralConfig()
	InitRPC()
}

func main() {
	probingCmd, err := QueryTask()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(probingCmd)

	rawData := Probe(probingCmd)
	jsonParams := Parse(rawData)
	Push(jsonParams)
}
