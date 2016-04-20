package main

import "log"

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
	log.Println("Execution the probing command:", probingCmd[0])

	rawData := Probe(probingCmd)
	jsonParams := MarshalIntoParameters(rawData)
	Push(jsonParams)
}
