package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Cepave/common/model"
	"github.com/montanaflynn/stats"
)

type Utility interface {
	CalcStats(row []float64, length int) map[string]string
	marshalStatsIntoJsonParams(stats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent
	ProbingCommand(targetAddressList []string) []string
	utilName() string
}

type Fping struct {
	Utility
}

func (fping *Fping) marshalStatsIntoJsonParams(fpingStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
	var params []ParamToAgent

	for rowNum, fpingStat := range fpingStats {
		// To graph
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-sent", fpingStat["pkttransmit"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-received", fpingStat["pktreceive"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "transmission-time", fpingStat["rttavg"]))

		// To Cassandra
		t := targetToNqmEndpointData(&targets[rowNum])
		a := agentToNqmEndpointData(agentPtr)
		nqmDataGram := TagsAssembler(t, a, fpingStat)
		params = append(params, nqmMarshalJSON(nqmDataGram, "nqm-fping"))
	}
	return params
}

func (fping *Fping) ProbingCommand(targetAddressList []string) []string {
	probingCmd := append(GetGeneralConfig().hbsResp.Command, targetAddressList...)
	return probingCmd
}

func (fping *Fping) utilName() string {
	return "fping"
}

func (fping *Fping) CalcStats(row []float64, length int) map[string]string {
	dataMap := map[string]string{
		"rttmin":    "-1",
		"rttmax":    "-1",
		"rttavg":    "-1",
		"rttmdev":   "-1",
		"rttmedian": "-1",
	}

	pktxmt := length
	pktrcv := len(row)
	var d stats.Float64Data = row
	median, _ := d.Median()
	max, _ := d.Max()
	min, _ := d.Min()
	mean, _ := d.Mean()
	dev, _ := d.StandardDeviation()

	if len(row) > 0 {
		dataMap["rttmin"] = strconv.FormatFloat(min, 'f', 2, 64)
		dataMap["rttmax"] = strconv.FormatFloat(max, 'f', 2, 64)
		dataMap["rttavg"] = strconv.FormatFloat(mean, 'f', 2, 64)
		dataMap["rttmdev"] = strconv.FormatFloat(dev, 'f', 2, 64)
		dataMap["rttmedian"] = strconv.FormatFloat(median, 'f', 2, 64)
	}
	dataMap["pkttransmit"] = strconv.Itoa(pktxmt)
	dataMap["pktreceive"] = strconv.Itoa(pktrcv)

	return dataMap
}

type Tcpping struct {
	Utility
}

func (tcpping *Tcpping) marshalStatsIntoJsonParams(tcppingStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
	return nil
}

func (tcpping *Tcpping) ProbingCommand(targetAddressList []string) []string {
	probingCmd := append([]string{"tcpping"}, targetAddressList...)
	return probingCmd
}

func (tcpping *Tcpping) utilName() string {
	return "tcpping"
}

type Tcpconn struct {
	Utility
}

func (tcpconn *Tcpconn) marshalStatsIntoJsonParams(tcpconnStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
	var params []ParamToAgent

	for rowNum, tcpconnStat := range tcpconnStats {
		// To graph
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "time", tcpconnStat["time"]))

		// To Cassandra
		t := targetToNqmEndpointData(&targets[rowNum])
		a := agentToNqmEndpointData(agentPtr)
		nqmDataGram := TagsAssembler(t, a, tcpconnStat)
		params = append(params, nqmMarshalJSON(nqmDataGram, "nqm-tcpconn"))
	}
	return params
}

func (tcpconn *Tcpconn) ProbingCommand(targetAddressList []string) []string {
	probingCmd := append([]string{"tcpconn"}, targetAddressList...)
	probingCmd = append(probingCmd, "| awk '$5{print $2\" : \"$5} !$5{print $2\" : -\"}'")
	probingCmd = []string{"/bin/sh", "-c", strings.Join(probingCmd, " ")}
	return probingCmd
}

func (tcpconn *Tcpconn) utilName() string {
	return "tcpconn"
}

func getTargetAddressList() []string {
	var targetAddressList []string
	for _, target := range GetGeneralConfig().hbsResp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}
	return targetAddressList
}

func (tcpconn *Tcpconn) CalcStats(row []float64, length int) map[string]string {
	dataMap := map[string]string{
		"time": "-1",
	}
	if length != 1 {
		log.Fatalln("Calculate statistic of tcpconn error")
	}
	if len(row) > 0 {
		time := row[0]
		dataMap["time"] = strconv.FormatFloat(time, 'f', 2, 64)
	}
	return dataMap
}

func measureByUtil(u Utility) {
	for {
		func() {
			probingCmd, err := Task(u)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("[", u.utilName(), "] Measuring...")

			rawData := Probe(probingCmd, u.utilName())
			parsedData := Parse(rawData)
			statsData := Calc(parsedData, u)
			jsonParams := u.marshalStatsIntoJsonParams(statsData, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent)

			for i, _ := range jsonParams {
				println(jsonParams[i].String())
				println("===")
			}

			Push(jsonParams)
		}()

		dur := time.Second * time.Duration(GetGeneralConfig().Agent.FpingInterval)
		time.Sleep(dur)
	}
}

func Measure() {
	go measureByUtil(new(Fping))
	//go measureByUtil(new(Tcpping))
	go measureByUtil(new(Tcpconn))
}
