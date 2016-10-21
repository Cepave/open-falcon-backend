package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
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

type nqmNodeData struct {
	Id         string
	IspId      string
	ProvinceId string
	CityId     string
	NameTagId  string
}

func marshalJSONParamsToCassandra(nqmDataGram string, metric string) ParamToAgent {
	data := ParamToAgent{}
	data.Tags = nqmDataGram
	data.Metric = metric
	data.Timestamp = time.Now().Unix()
	data.Endpoint = Meta().Hostname
	data.Value = "0"
	data.CounterType = "GAUGE"
	data.Step = int64(60) // a useless field in Cassandra
	return data
}

func convToKeyValueString(arg map[string]string) string {
	Str := ""
	for key, value := range arg {
		Str = Str + "," + key + "=" + value
	}
	return Str
}

func assembleTags(target nqmNodeData, agent nqmNodeData, dataMap map[string]string) string {
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
		convToKeyValueString(dataMap)
}

func convToNqmAgent(s model.NqmAgent) nqmNodeData {
	return nqmNodeData{
		Id:         strconv.Itoa(s.Id),
		IspId:      strconv.Itoa(int(s.IspId)),
		ProvinceId: strconv.Itoa(int(s.ProvinceId)),
		CityId:     strconv.Itoa(int(s.CityId)),
		NameTagId:  strconv.Itoa(-1),
	}
}

func convToNqmTarget(s model.NqmTarget) nqmNodeData {
	return nqmNodeData{
		Id:         strconv.Itoa(s.Id),
		IspId:      strconv.Itoa(int(s.IspId)),
		ProvinceId: strconv.Itoa(int(s.ProvinceId)),
		CityId:     strconv.Itoa(int(s.CityId)),
		NameTagId:  strconv.Itoa(-1),
	}
}

/**
 * value could be:
 *     Packet Loss - int
 *     Transmission Time - float64
 */
func marshalJSONToGraph(target model.NqmTarget, agent model.NqmAgent, metric string, value interface{}, step int64) ParamToAgent {
	endpoint := Meta().Hostname
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
	return ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
}

func marshalStatsRow(row map[string]string, target model.NqmTarget, agent model.NqmAgent, step int64, u Utility) []ParamToAgent {
	var params []ParamToAgent

	params = append(params, u.MarshalJSONParamsToGraph(target, agent, row, step)...)

	t := convToNqmTarget(target)
	a := convToNqmAgent(agent)
	cassandraDataGram := assembleTags(t, a, row)
	params = append(params, marshalJSONParamsToCassandra(cassandraDataGram, "nqm-"+u.UtilName()))

	return params
}

func Marshal(statsData []map[string]string, u Utility, targets []model.NqmTarget, agent model.NqmAgent, step int64) []ParamToAgent {
	var params []ParamToAgent

	for rowIdx, statsRow := range statsData {
		target := targets[rowIdx]
		params = append(params, marshalStatsRow(statsRow, target, agent, step, u)...)
	}
	return params
}
