package main

import "log"

//*
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	InitGeneralConfig()
	InitRPC()
}

//*/
func main() {
	probingCmd, targets, err := QueryTask()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Execution the probing command:", probingCmd[0])

	rawData := Probe(probingCmd)
	jsonParams := MarshalIntoParameters(rawData, targets)
	Push(jsonParams)
	//	log.Println(jsonParams)
}
