package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdateMeasurements(t *testing.T) {
	var cfg GeneralConfig
	generalConfig = &cfg
	cfg.Agent = new(AgentConfig)
	cfg.Agent.FpingInterval = 60
	cfg.Agent.TcppingInterval = 60
	cfg.Agent.TcpconnInterval = 60

	tests := [][]string{
		{"", "", ""},
		{"", "", "tcpconn"},
		{"", "tcpping", ""},
		{"", "tcpping", "tcpconn"},
		{"fping", "", ""},
		{"fping", "", "tcpconn"},
		{"fping", "tcpping", ""},
		{"fping", "tcpping", "tcpconn"},
	}

	expecteds := []map[string]MeasurementsProperty{
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, false}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, true}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, true}},
	}

	for j, c := range tests {
		d := updateMeasurements(c)
		if !reflect.DeepEqual(expecteds[j], d) {
			t.Error(expecteds[j], d)
		}
		t.Log(expecteds[j], d)

	}
}

func TestUpdatedMsg(t *testing.T) {
	testsOld := []map[string]MeasurementsProperty{
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, false}},
	}

	testsUpdated := []map[string]MeasurementsProperty{
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, false}},
	}

	expecteds := []string{
		"<fping Disabled> <tcpping Disabled> <tcpconn Disabled> ",
		"<fping Enabled> <tcpping Enabled> <tcpconn Enabled> ",
		"<fping Disabled> ",
		"<fping Disabled> <tcpconn Enabled> ",
		"",
	}

	for i, _ := range testsOld {
		o := testsOld[i]
		u := testsUpdated[i]
		msg := updatedMsg(o, u)

		msgSlice := strings.SplitAfter(msg, "> ")
		msgMap := make(map[string]bool)
		for _, m := range msgSlice {
			if m == "" {
				continue
			}
			msgMap[m] = true
		}

		expectedSlice := strings.SplitAfter(expecteds[i], "> ")
		expectedMap := make(map[string]bool)
		for _, m := range expectedSlice {
			if m == "" {
				continue
			}
			expectedMap[m] = true
		}

		if !reflect.DeepEqual(expectedMap, msgMap) {
			t.Error(expectedMap, msgMap)
		}
		t.Log(msg)
	}
}
