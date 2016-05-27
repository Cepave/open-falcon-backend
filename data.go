package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Cepave/common/model"
)

type ParamToAgent struct {
	Metric      string      `json:"metric"`
	Endpoint    string      `json:"endpoint"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	Tags        string      `json:"tags"`
	Timestamp   int64       `json:"timestamp"`
	Step        int64       `json:"step"`
}

func (p ParamToAgent) String() string {
	return fmt.Sprintf(
		" {metric: %v, endpoint: %v, value: %v, counterType:%v, tags:%v, timestamp:%d, step:%d}",
		p.Metric,
		p.Endpoint,
		p.Value,
		p.CounterType,
		p.Tags,
		p.Timestamp,
		p.Step,
	)
}

type nqmEndpointData struct {
	Id         string
	IspId      string
	ProvinceId string
	CityId     string
	NameTagId  string
}

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

func TagsAssembler(target *nqmEndpointData, agent *nqmEndpointData, nqmDataMap map[string]string) string {
	return "agent-id=" + agent.Id +
		",agent-isp-id=" + agent.IspId +
		",agent-province-id=" + agent.ProvinceId +
		",agent-city-id=" + agent.CityId +
		",agent-name-tag-id=" + agent.NameTagId +
		",target-id=" + target.Id +
		",target-isp-id=" + target.IspId +
		",target-province-id=" + target.ProvinceId +
		",target-city-id=" + target.CityId +
		",target-name-tag-id=" + target.NameTagId +
		",rttmin=" + nqmDataMap["rttmin"] +
		",rttmax=" + nqmDataMap["rttmax"] +
		",rttavg=" + nqmDataMap["rttavg"] +
		",rttmdev=" + nqmDataMap["rttmdev"] +
		",rttmedian=" + nqmDataMap["rttmedian"] +
		",pkttransmit=" + nqmDataMap["pkttransmit"] +
		",pktreceive=" + nqmDataMap["pktreceive"] +
		",time=" + nqmDataMap["time"]
}

func nqmMarshalJSON(nqmDataGram string, metric string) ParamToAgent {
	data := ParamToAgent{}
	data.Tags = nqmDataGram
	data.Metric = metric
	data.Timestamp = time.Now().Unix()
	data.Endpoint = GetGeneralConfig().Hostname
	data.Value = "0"
	data.CounterType = "GAUGE"
	data.Step = int64(60)
	return data
}

/**
 * value could be:
 *     Packet Loss - int
 *     Transmission Time - float64
 */

func marshalJSON(target model.NqmTarget, agent *model.NqmAgent, metric string, value interface{}) ParamToAgent {
	endpoint := GetGeneralConfig().Hostname
	counterType := "GAUGE"
	tags := "nqm-agent-isp=" + agent.IspName +
		",nqm-agent-province=" + agent.ProvinceName +
		",nqm-agent-city=" + agent.CityName +
		",target-ip=" + target.Host +
		",target-isp=" + target.IspName +
		",target-province=" + target.ProvinceName +
		",target-city=" + target.CityName +
		",target-name-tag=" + target.NameTag
	timestamp := time.Now().Unix()
	step := int64(60)
	return ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
}
