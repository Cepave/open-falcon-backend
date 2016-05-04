package main

import (
	"fmt"
	"github.com/Cepave/common/model"
	"log"
)

func QueryTask() ([]string, []model.NqmTarget, *model.NqmAgent, error) {
	err := rpcClient.Call("NqmAgent.PingTask", req, &resp)
	if err != nil {
		log.Fatalln("Call NqmAgent.PingTask error:", err)
	}
	if !resp.NeedPing {
		return []string{}, resp.Targets, resp.Agent, fmt.Errorf("No tasks assigned.")
	}

	var targetAddressList []string
	for _, target := range resp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append(resp.Command, targetAddressList...)
	return probingCmd, resp.Targets, resp.Agent, err
}
