package main

import (
	"log"
	"strings"
	"time"

	"github.com/Cepave/common/model"
)

type Utility interface {
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
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-sent", fpingStat["pkttransmit"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-received", fpingStat["pktreceive"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "transmission-time", fpingStat["rttavg"]))

		t := targetToNqmEndpointData(&targets[rowNum])
		a := agentToNqmEndpointData(agentPtr)
		nqmDataGram := nqmTagsAssembler(t, a, fpingStat)
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
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "tcpconntime", tcpconnStat["tcpconntime"]))
		t := targetToNqmEndpointData(&targets[rowNum])
		a := agentToNqmEndpointData(agentPtr)
		nqmDataGram := nqmTagsAssembler(t, a, tcpconnStat)
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

func statsCalc(parsedData [][]string) []map[string]string {
	var stats []map[string]string
	for _, row := range parsedData {
		stat := nqmFpingStat(row, "fping")
		stats = append(stats, stat)
	}
	return stats
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
			stats := statsCalc(parsedData)
			jsonParams := u.marshalStatsIntoJsonParams(stats, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent)

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
