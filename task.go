package main

import (
	"fmt"
	"github.com/Cepave/common/model"
	"log"
	"strconv"
)

type nqmEndpointData struct {
	Id         string
	IspId      string
	ProvinceId string
	CityId     string
	NameTagId  string
}

var (
	agentData *nqmEndpointData
)

func agentToNqmEndpointData(s *model.NqmAgent) *nqmEndpointData {
	return &nqmEndpointData{
		Id:         strconv.Itoa(s.Id),
		IspId:      strconv.Itoa(int(s.IspId)),
		ProvinceId: strconv.Itoa(int(s.ProvinceId)),
		CityId:     strconv.Itoa(int(s.CityId)),
		NameTagId:  strconv.Itoa(-1),
	}
}

func targetToNqmEndpointData(s *model.NqmTarget) *nqmEndpointData {
	return &nqmEndpointData{
		Id:         strconv.Itoa(s.Id),
		IspId:      strconv.Itoa(int(s.IspId)),
		ProvinceId: strconv.Itoa(int(s.ProvinceId)),
		CityId:     strconv.Itoa(int(s.CityId)),
		NameTagId:  strconv.Itoa(-1),
	}
}

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
	agentData = agentToNqmEndpointData(&agent)

	var targetAddressList []string
	for _, target := range resp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append(resp.Command, targetAddressList...)
	return probingCmd, nil
}
