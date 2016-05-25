package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Cepave/common/model"
)

type Utility interface {
	marshalStatsIntoJSON([]map[string]string, []model.NqmTarget, *model.NqmAgent) []ParamToAgent
	task() ([]string, []model.NqmTarget, *model.NqmAgent, error)
	utilName() string
}

type Fping struct {
	Utility
}

func (fping *Fping) marshalStatsIntoJSON(fpingStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
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

func (fping *Fping) task() ([]string, []model.NqmTarget, *model.NqmAgent, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, nil, nil, fmt.Errorf("[ " + fping.utilName() + " ] No tasks assigned.")
	}

	var targetAddressList []string
	for _, target := range GetGeneralConfig().hbsResp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append(GetGeneralConfig().hbsResp.Command, targetAddressList...)
	return probingCmd, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent, nil
}

func (fping *Fping) utilName() string {
	return "fping"
}

/*
type tcppingMeasure struct {
	measure
}

func (tcpping *tcppingMeasure) probe(probingCmd []string) []string {
	cmdOutput, err := exec.Command(probingCmd[0], probingCmd[1:]...).CombinedOutput()
	if err != nil {
		// 'exit status 1' happens when there is at least
		// one target with 100% packet loss.
		log.Println("[ tcpping ] An error occured:", err)
	}
	tcppingResults := strings.Split(string(cmdOutput), "\n")
	tcppingResults = tcppingResults[:len(tcppingResults)-1]
	rawData := trimResults(tcppingResults)
	return rawData
}

func (tcpping *tcppingMeasure) marshalStatsIntoJSON(tcppingStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
	var params []ParamToAgent

	for rowNum, tcppingStat := range tcppingStats {
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-sent", tcppingStat["pkttransmit"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "packets-received", tcppingStat["pktreceive"]))
		params = append(params, marshalJSON(targets[rowNum], agentPtr, "transmission-time", tcppingStat["rttavg"]))

		t := targetToNqmEndpointData(&targets[rowNum])
		a := agentToNqmEndpointData(agentPtr)
		nqmDataGram := nqmTagsAssembler(t, a, tcppingStat)
		params = append(params, nqmMarshalJSON(nqmDataGram, "nqm-tcpping"))
	}
	return params
}
*/

type Tcpconn struct {
	Utility
}

func (tcpconn *Tcpconn) marshalStatsIntoJSON(tcpconnStats []map[string]string, targets []model.NqmTarget, agentPtr *model.NqmAgent) []ParamToAgent {
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

func (tcpconn *Tcpconn) task() ([]string, []model.NqmTarget, *model.NqmAgent, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, nil, nil, fmt.Errorf("[ " + tcpconn.utilName() + " ] No tasks assigned.")
	}

	var targetAddressList []string
	for _, target := range GetGeneralConfig().hbsResp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}

	probingCmd := append([]string{"/home/vagrant/workspace/shell/nqm/tcpconn/tcpconn"}, targetAddressList...)
	probingCmd = append(probingCmd, "| awk '$5{print $2\" : \"$5} !$5{print $2\" : -\"}'")
	probingCmd = []string{"/bin/sh", "-c", strings.Join(probingCmd, " ")}
	return probingCmd, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent, nil
}

func (tcpconn *Tcpconn) utilName() string {
	return "tcpconn"
}

func statsCalc(parsedData [][]string) []map[string]string {
	var stats []map[string]string
	for _, row := range parsedData {
		stat := nqmFpingStat(row, "fping")
		stats = append(stats, stat)
	}
	return stats
}

func parse(rawData []string) [][]string {
	var parsedRows [][]string
	for _, row := range rawData {
		parsedRow := parseFpingRow(row)
		parsedRows = append(parsedRows, parsedRow)
	}
	return parsedRows
}

func measureBy(u Utility) {
	for {
		func() {
			probingCmd, targets, agentPtr, err := u.task()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("[", u.utilName(), "] Measuring...")

			rawData := Probe(probingCmd, u.utilName())
			for _, row := range rawData {
				fmt.Println(row)
			}
			parsedData := parse(rawData)
			fpingStats := statsCalc(parsedData)
			jsonParams := u.marshalStatsIntoJSON(fpingStats, targets, agentPtr)

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
	//go fping()
	//go tcpping()
	//go tcpconn()
	//go fpingRefactor()
	go measureBy(new(Fping))
	go measureBy(new(Tcpconn))
}
