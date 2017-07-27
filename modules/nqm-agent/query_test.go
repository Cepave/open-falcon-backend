package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Cepave/open-falcon-backend/common/model"
)

func TestUpdatedMsg(t *testing.T) {
	testsOld := []map[string]model.MeasurementsProperty{
		{"fping": {true, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
		{"fping": {true, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		{"fping": {true, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		nil,
	}

	testsUpdated := []map[string]model.MeasurementsProperty{
		{"fping": {false, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
		{"fping": {true, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {false, []string{""}, 60}, "tcpconn": {true, []string{""}, 60}},
		{"fping": {false, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
		nil,
		{"fping": {true, []string{""}, 60}, "tcpping": {true, []string{""}, 60}, "tcpconn": {false, []string{""}, 60}},
	}

	expecteds := []string{
		"<fping Disabled> <tcpping Disabled> <tcpconn Disabled> ",
		"<fping Enabled> <tcpping Enabled> <tcpconn Enabled> ",
		"<fping Disabled> ",
		"<fping Disabled> <tcpconn Enabled> ",
		"",
		"<tcpping Disabled> <tcpconn Disabled> ",
		"<fping Enabled> <tcpping Enabled> ",
	}

	for i := range testsOld {
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
