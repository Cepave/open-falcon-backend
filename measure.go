package main

import (
	"log"
	"strconv"
	"time"

	"github.com/Cepave/common/model"
	"github.com/montanaflynn/stats"
)

type Utility interface {
	CalcStats(row []float64, length int) map[string]string
	MarshalJSONParamsToGraph(target model.NqmTarget, agent model.NqmAgent, row map[string]string) []ParamToAgent
	ProbingCommand(command []string, targetAddressList []string) []string
	UtilName() string
	Interval() time.Duration
}

type Fping struct {
	Utility
}

func (u *Fping) MarshalJSONParamsToGraph(target model.NqmTarget, agent model.NqmAgent, row map[string]string) []ParamToAgent {
	var params []ParamToAgent

	params = append(params, marshalJSONToGraph(target, agent, "packets-sent", row["pkttransmit"]))
	params = append(params, marshalJSONToGraph(target, agent, "packets-received", row["pktreceive"]))
	params = append(params, marshalJSONToGraph(target, agent, "transmission-time", row["rttavg"]))

	return params
}

func (u *Fping) ProbingCommand(command []string, targetAddressList []string) []string {
	probingCmd := append(command, targetAddressList...)
	return probingCmd
}

func (u *Fping) UtilName() string {
	return "fping"
}

func (u *Fping) CalcStats(row []float64, length int) map[string]string {
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

func (u *Fping) Interval() time.Duration {
	return time.Second * time.Duration(GetGeneralConfig().Agent.FpingInterval)
}

type Tcpping struct {
	Utility
}

func (u *Tcpping) MarshalJSONParamsToGraph(target model.NqmTarget, agent model.NqmAgent, row map[string]string) []ParamToAgent {
	return new(Fping).MarshalJSONParamsToGraph(target, agent, row)
}

func (u *Tcpping) ProbingCommand(command []string, targetAddressList []string) []string {
	probingCmd := append([]string{"tcpping.sh"}, targetAddressList...)
	return probingCmd
}

func (u *Tcpping) UtilName() string {
	return "tcpping"
}

func (u *Tcpping) CalcStats(row []float64, length int) map[string]string {
	return new(Fping).CalcStats(row, length)
}

func (u *Tcpping) Interval() time.Duration {
	return time.Second * time.Duration(GetGeneralConfig().Agent.TcppingInterval)
}

type Tcpconn struct {
	Utility
}

func (u *Tcpconn) MarshalJSONParamsToGraph(target model.NqmTarget, agent model.NqmAgent, row map[string]string) []ParamToAgent {
	var params []ParamToAgent
	params = append(params, marshalJSONToGraph(target, agent, "tcpconn", row["time"]))
	return params
}

func (u *Tcpconn) ProbingCommand(command []string, targetAddressList []string) []string {
	probingCmd := append([]string{"tcpconn.sh"}, targetAddressList...)
	return probingCmd
}

func (u *Tcpconn) UtilName() string {
	return "tcpconn"
}

func (u *Tcpconn) CalcStats(row []float64, length int) map[string]string {
	dataMap := map[string]string{
		"time": "-1",
	}
	if length != 1 {
		log.Fatalln("[", u.UtilName(), "] Error on Calculation: Wrong format of the statistics")
	}
	if len(row) > 0 {
		time := row[0]
		dataMap["time"] = strconv.FormatFloat(time, 'f', 2, 64)
	}
	return dataMap
}

func (u *Tcpconn) Interval() time.Duration {
	return time.Second * time.Duration(GetGeneralConfig().Agent.TcpconnInterval)
}

func measureByUtil(u Utility) {
	probingCmd, targets, agent, err := Task(u)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("[", u.UtilName(), "] Measuring...")

	rawData := Probe(probingCmd, u.UtilName())
	parsedData := Parse(rawData)
	statsData := Calc(parsedData, u)
	jsonParams := Marshal(statsData, u, targets, agent)
	Push(jsonParams, u.UtilName())
}

func measure(u Utility) {
	for {
		go measureByUtil(u)
		time.Sleep(u.Interval())
	}
}

func Measure() {
	go measure(new(Fping))
	go measure(new(Tcpping))
	go measure(new(Tcpconn))
}
