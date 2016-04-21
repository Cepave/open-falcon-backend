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

func (p ParamToAgent) String() string {
	return fmt.Sprintf(
		" {metric: %v, endpoint: %v, value: %v, counterType:%v, tags:%v, timestamp:%d, step:%d}",
		p.Metric,
		p.Endpoint,
		p.Value,
		p.CounterType,
		p.Tags,
		p.Timestamp,
		p.Step,
	)
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
	log.Println("Execution the probing command:", probingCmd[0])

	rawData := Probe(probingCmd)
	jsonParams := MarshalIntoParameters(rawData)
	Push(jsonParams)
	//	log.Println(jsonParams)
}
