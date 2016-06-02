package main

import (
	"reflect"
	"testing"
)

func TestUpdateMeasurements(t *testing.T) {
	var cfg GeneralConfig
	generalConfig = &cfg
	cfg.Agent = new(AgentConfig)
	cfg.Agent.FpingInterval = 60
	cfg.Agent.TcppingInterval = 60
	cfg.Agent.TcpconnInterval = 60

	testsMeasurements := []map[string]MeasurementsProperty{
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, false}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, false}},
		{"fping": {60, false}, "tcpping": {60, true}, "tcpconn": {60, true}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, false}, "tcpconn": {60, true}},
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, false}},
		{"fping": {60, true}, "tcpping": {60, true}, "tcpconn": {60, true}},
	}

	testsCmd := [][]string{
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

	for _, m := range testsMeasurements {
		for j, c := range testsCmd {
			d := updateMeasurements(m, c)
			if !reflect.DeepEqual(expecteds[j], d) {
				t.Error(expecteds[j], d)
			}
			t.Log(expecteds[j], d)
		}
	}
}
