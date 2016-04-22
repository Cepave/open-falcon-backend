package main

import (
	"github.com/Cepave/common/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func parseFpingRow(row string) []string {
	return strings.FieldsFunc(row, func(r rune) bool {
		switch r {
		case ' ', '\n', ':', '/', '%', '=', ',':
			return true
		}
		return false
	})
}

func marshalFpingRowIntoJSON(row []string, target model.NqmTarget) []ParamToAgent {
	var params []ParamToAgent
	xmt, err := strconv.Atoi(row[4])
	if err != nil {
		log.Println("error occured:", err)
	}
	params = append(params, marshalJSON(target, "packets-sent", xmt))

	rcv, err := strconv.Atoi(row[5])
	if err != nil {
		log.Println("error occured:", err)
	}
	params = append(params, marshalJSON(target, "packets-received", rcv))

	tt, err := strconv.ParseFloat(row[11], 64)
	if err != nil {
		log.Println("error occured:", err)
	}
	params = append(params, marshalJSON(target, "transmission-time", tt))

	params = append(params, nqmMarshalJSON(target, "nqm-metrics", row))
	return params
}

func nqmParseFpingRow(row []string) map[string]string {
	/*
		www.yahoo.com  : xmt/rcv/%loss = 100/99/1%, min/avg/max = 5.42/10.9/35.9
		                                  4  5  6                  10   11  12
	*/
	nqmDataMap := map[string]string{}
	nqmDataMap["rttmin"] = row[10]
	nqmDataMap["rttmax"] = row[12]
	nqmDataMap["rttavg"] = row[11]
	nqmDataMap["rttmdev"] = "-1"
	nqmDataMap["rttmedian"] = "-1"
	nqmDataMap["pkttransmit"] = row[4]
	nqmDataMap["pktreceive"] = row[5]
	return nqmDataMap
}

func nqmTagsAssembler(target *nqmEndpointData, agent *nqmEndpointData, nqmDataMap map[string]string) string {
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
		",pktreceive=" + nqmDataMap["pktreceive"]
}

func nqmMarshalJSON(target model.NqmTarget, metric string, row []string) ParamToAgent {
	t := targetToNqmEndpointData(&target)
	data := ParamToAgent{}
	data.Tags = nqmTagsAssembler(t, agentData, nqmParseFpingRow(row))
	data.Metric = metric
	data.Timestamp = time.Now().Unix()
	data.Endpoint = "nqm-endpoint"
	data.Value = "0"
	data.CounterType = "nqm"
	return data
}

/**
 * value could be:
 *     Packet Loss - int
 *     Transmission Time - float64
 */
func marshalJSON(target model.NqmTarget, metric string, value interface{}) ParamToAgent {
	endpoint := GetGeneralConfig().Hostname
	counterType := "GAUGE"
	tags := "nqm-agent-isp=" + GetGeneralConfig().ISP +
		",nqm-agent-province=" + GetGeneralConfig().Province +
		",nqm-agent-city=" + GetGeneralConfig().City +
		",target-ip=" + target.Host +
		",target-isp=" + target.IspName +
		",target-province=" + target.ProvinceName +
		",target-city=" + target.CityName +
		",target-name-tag=" + target.NameTag
	timestamp := time.Now().Unix()
	step := int64(60)
	return ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
}

func MarshalIntoParameters(rawData []string) []ParamToAgent {
	var params []ParamToAgent
	for rowNum, row := range rawData {
		parsedRow := parseFpingRow(row)
		if len(parsedRow) != 13 {
			continue
		}

		target := resp.Targets[rowNum]
		params = append(params, marshalFpingRowIntoJSON(parsedRow, target)...)
	}
	return params
}
