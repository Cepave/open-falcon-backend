package main

import (
	"log"
	"reflect"
	"time"

	"github.com/Cepave/common/model"
)

func configFromHbsUpdated(newResp model.NqmPingTaskResponse) bool {
	if !reflect.DeepEqual(GetGeneralConfig().hbsResp, newResp) {
		return true
	}
	return false
}

func query() {
	var resp model.NqmPingTaskResponse
	err := rpcClient.Call("NqmAgent.PingTask", req, &resp)
	if err != nil {
		log.Fatalln("Call NqmAgent.PingTask error:", err)
	}
	log.Println("[ hbs ] Response received")

	if configFromHbsUpdated(resp) {
		GetGeneralConfig().hbsResp = resp
		log.Println("[ hbs ] Configuration updated")
	}
}

func Query() {
	for {
		go query()

		dur := time.Second * time.Duration(GetGeneralConfig().Hbs.Interval)
		time.Sleep(dur)
	}
}
