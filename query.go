package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Cepave/common/model"
)

func updatedMsg(old map[string]MeasurementsProperty, updated map[string]MeasurementsProperty) string {
	msg := ""
	for k, _ := range updated {
		if !old[k].enabled && updated[k].enabled {
			msg = msg + fmt.Sprint("<", k, " Enabled> ")
		}
		if old[k].enabled && !updated[k].enabled {
			msg = msg + fmt.Sprint("<", k, " Disabled> ")
		}
	}
	return msg
}

func updateMeasurements(command []string) map[string]MeasurementsProperty {
	updated := NewMeasurements()

	for _, cmd := range command {
		if m, ok := updated[cmd]; ok {
			m.enabled = true
			updated[cmd] = m
		}
	}
	return updated
}

func configFromHbsUpdated(newResp model.NqmPingTaskResponse) bool {
	if !reflect.DeepEqual(GetGeneralConfig().hbsResp.Load().(model.NqmPingTaskResponse), newResp) {
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

	GetGeneralConfig().hbsResp.Store(resp)

	old := GetGeneralConfig().Measurements
	updated := updateMeasurements(resp.Command)
	GetGeneralConfig().Measurements = updated

	if msg := updatedMsg(old, updated); msg != "" {
		log.Println("[ hbs ]", msg)
	}
	log.Println("[ hbs ] Configuration updated")
}

func Query() {
	for {
		query()

		dur := time.Second * time.Duration(GetGeneralConfig().Hbs.Interval)
		time.Sleep(dur)
	}
}
