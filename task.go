package main

import (
	"fmt"
	"log"
)

func QueryTask() ([]string, error) {
	err := rpcClient.Call("NqmAgent.PingTask", req, &resp)
	if err != nil {
		log.Fatalln("Call NqmAgent.PingTask error:", err)
	}
	if !resp.NeedPing {
		return []string{}, fmt.Errorf("No tasks assigned.")
	}
	agent := *resp.Agent
	SetGeneralConfigByAgent(agent)

	var targetAddressList []string
	for _, target := range resp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append(resp.Command, targetAddressList...)
	return probingCmd, nil
}
