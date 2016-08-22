package main

import (
	"fmt"
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
)

func tick() <-chan time.Time {
	dur := time.Second * GetGeneralConfig().Hbs.Interval
	return time.Tick(dur)
}

func updatedMsg(old map[string]model.MeasurementsProperty, updated map[string]model.MeasurementsProperty) string {
	msg := ""
	if updated == nil { // Reset all enabled measurements
		for k, v := range old {
			if v.Enabled {
				msg = msg + fmt.Sprint("<", k, " Disabled> ")
			}
		}
		return msg
	}
	for k, _ := range updated {
		if !old[k].Enabled && updated[k].Enabled {
			msg = msg + fmt.Sprint("<", k, " Enabled> ")
		}
		if old[k].Enabled && !updated[k].Enabled {
			msg = msg + fmt.Sprint("<", k, " Disabled> ")
		}
	}
	return msg
}

func configFromHbsUpdated(newResp model.NqmTaskResponse, oldResp model.NqmTaskResponse) bool {
	if !reflect.DeepEqual(newResp, oldResp) {
		return true
	}
	return false
}

func query() {
	var resp model.NqmTaskResponse
	err := RPCCall("NqmAgent.Task", req, &resp)
	if err != nil {
		log.Println("[ hbs ] Error on RPC call:", err)
		return
	}
	log.Println("[ hbs ] Response received")
	HbsRespTime = time.Now()

	oldResp := GetGeneralConfig().hbsResp.Load().(model.NqmTaskResponse)
	if !configFromHbsUpdated(resp, oldResp) {
		return
	}
	msg := "[ hbs ] Configuration updated"

	oldMeas := oldResp.Measurements
	updatedMeas := resp.Measurements
	if measMsg := updatedMsg(oldMeas, updatedMeas); measMsg != "" {
		msg = msg + " - " + measMsg
	}

	GetGeneralConfig().hbsResp.Store(resp)
	log.Println(msg)
}

func Query() {
	query()
	for _ = range tick() {
		query()
	}
}
