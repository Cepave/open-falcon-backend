package main

import (
	"fmt"
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
	log.Println("A response from hbs.")

	if configFromHbsUpdated(resp) {
		GetGeneralConfig().hbsResp = resp
		log.Println("The configuration was updated by hbs.")
	}
}

func makeTasks() ([]string, []model.NqmTarget, *model.NqmAgent, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, nil, nil, fmt.Errorf("No tasks assigned.")
	}

	var targetAddressList []string
	for _, target := range GetGeneralConfig().hbsResp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append(GetGeneralConfig().hbsResp.Command, targetAddressList...)
	return probingCmd, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent, nil
}

func QueryHbs() {
	for {
		query()

		dur := time.Second * time.Duration(GetGeneralConfig().Hbs.Interval)
		time.Sleep(dur)
	}
}
