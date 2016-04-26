package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Cepave/common/model"
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

	return params
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
