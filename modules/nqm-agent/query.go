package main

import (
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
)

var hbsTicker *time.Ticker
var hbsTickerUpdated chan bool

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
	for k := range updated {
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
		log.Errorln("[ hbs ] Error on RPC call:", err)
		return
	}
	log.Println("[ hbs ] Response received")
	HbsRespTime = time.Now()

	oldResp := HBSResp()
	if !configFromHbsUpdated(resp, oldResp) {
		return
	}
	msg := "[ hbs ] Configuration updated"

	oldMeas := oldResp.Measurements
	updatedMeas := resp.Measurements
	if measMsg := updatedMsg(oldMeas, updatedMeas); measMsg != "" {
		msg = msg + " - " + measMsg
	}

	SetHBSResp(resp)
	log.Println(msg)
}

func Query() {
	for {
		select {
		case <-hbsTicker.C:
			query()
		case <-hbsTickerUpdated:
			hbsTicker.Stop()
			hbsTicker = time.NewTicker(Config().Hbs.Interval * time.Second)
		}
	}
}
