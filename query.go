package main

import (
	"log"
	"reflect"
	"time"

	"github.com/Cepave/common/model"
)

func updateMeasurements(curr map[string]MeasurementsProperty, command []string) map[string]MeasurementsProperty {
	updated := NewMeasurements()

	for _, cmd := range command {
		if m, ok := updated[cmd]; ok {
			m.enabled = true
			updated[cmd] = m
		}
	}

	for k, _ := range updated {
		if !curr[k].enabled && updated[k].enabled {
			log.Println("[ hbs ] Enable <", k, ">")
		}
		if curr[k].enabled && !updated[k].enabled {
			log.Println("[ hbs ] Disable <", k, ">")
		}
	}
	return updated
}

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
		log.Println("[ hbs ] Error on RPC call:", err)
		return
	}

	log.Println("[ hbs ] Response received")
	if !configFromHbsUpdated(resp) {
		return
	}

	GetGeneralConfig().hbsResp = resp
	curMeasurements := GetGeneralConfig().Measurements
	cmd := GetGeneralConfig().hbsResp.Command
	GetGeneralConfig().Measurements = updateMeasurements(curMeasurements, cmd)
	log.Println("[ hbs ] Configuration updated")

}

func Query() {
	for {
		go query()

		dur := time.Second * time.Duration(GetGeneralConfig().Hbs.Interval)
		time.Sleep(dur)
	}
}
