package main

import (
	"strings"
	"testing"
	"time"

	"github.com/Cepave/common/model"
)

func TestConvToKeyValueString(t *testing.T) {
	tests := []map[string]string{
		{"rttmax": "46.5", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8.61"},
		{"time": "13.6"},
	}

	expecteds := []string{
		",rttmin=8.61,rttmax=46.5,rttavg=14.5,rttmdev=-1,rttmedian=-1,pkttransmit=100,pktreceive=100",
		",time=13.6",
	}

	for i, v := range tests {

		kv := convToKeyValueString(v)
		if !strings.HasPrefix(kv, ",") {
			t.Error(kv)
		}

		strs := strings.Split(kv, ",")
		for _, str := range strs {
			if !strings.Contains(expecteds[i], str) {
				t.Error(expecteds[i], str)
			}
		}

		t.Log(convToKeyValueString(v))
	}
}

func TestAssembleTags(t *testing.T) {
	agent := nqmNodeData{
		"-1", "-1", "-1", "-1", "-1",
	}
	target := nqmNodeData{
		"-2", "-2", "-2", "-2", "-2",
	}
	tests := []map[string]string{
		{"rttmax": "46.5", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8.61"},
		{"time": "13.6"},
	}

	expecteds := []string{
		"agent-id=-1,agent-isp-id=-1,agent-province-id=-1,agent-city-id=-1,agent-name-tag-id=-1,target-id=-2,target-isp-id=-2,target-province-id=-2,target-city-id=-2,target-name-tag-id=-2,rttmin=8.61,rttmax=46.5,rttavg=14.5,rttmdev=-1,rttmedian=-1,pkttransmit=100,pktreceive=100,time=",
		"agent-id=-1,agent-isp-id=-1,agent-province-id=-1,agent-city-id=-1,agent-name-tag-id=-1,target-id=-2,target-isp-id=-2,target-province-id=-2,target-city-id=-2,target-name-tag-id=-2,rttmin=,rttmax=,rttavg=,rttmdev=,rttmedian=,pkttransmit=,pktreceive=,time=13.6",
	}

	//t_out := nqmTagsAssembler(target, agent, tests)
	for i, v := range tests {
		kv := assembleTags(target, agent, v)
		strs := strings.Split(kv, ",")
		for _, str := range strs {
			if !strings.Contains(expecteds[i], str) {
				t.Error(expecteds[i], str)
			}
		}
		t.Log(assembleTags(target, agent, v))
	}
}

func TestMarshalStatsRow(t *testing.T) {
	// Hostname is the config dependency which lies in func MarshalIntoParameters
	var cfg GeneralConfig
	generalConfig = &cfg
	cfg.Hostname = "unit-test-hostname"

	nqmAgent := model.NqmAgent{
		Id: -1, IspId: -1, ProvinceId: -1, CityId: -1,
	}
	nqmTarget := model.NqmTarget{
		Id: -2, IspId: -2, ProvinceId: -2, CityId: -2,
	}

	var resp model.NqmTaskResponse
	resp.NeedPing = true
	resp.Agent = &nqmAgent
	resp.Targets = []model.NqmTarget{nqmTarget}
	resp.Measurements = map[string]model.MeasurementsProperty{
		"fping":   {true, []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a"}, 300},
		"tcpping": {false, []string{"tcpping", "-i", "0.01", "-c", "4"}, 300},
		"tcpconn": {false, []string{"tcpconn"}, 300},
	}
	GetGeneralConfig().hbsResp.Store(resp)

	agent := nqmNodeData{
		"-1", "-1", "-1", "-1", "-1",
	}
	target := nqmNodeData{
		"-2", "-2", "-2", "-2", "-1",
	}

	// fping, 4 JSON parameters, only test the 4th which is for Cassandra
	tests := []map[string]string{
		{"rttmax": "45.08", "rttavg": "44.63", "rttmdev": "0.31", "rttmedian": "44.58", "pkttransmit": "4", "pktreceive": "4", "rttmin": "44.28"},
		{"rttmax": "-1", "rttavg": "-1", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "4", "pktreceive": "0", "rttmin": "-1"},
	}
	expecteds := [][]ParamToAgent{
		{{}, {}, {}, {Metric: "nqm-fping", Endpoint: "unit-test-hostname", Value: 0, CounterType: "GAUGE", Tags: "", Timestamp: time.Now().Unix(), Step: 60}},
		{{}, {}, {}, {Metric: "nqm-fping", Endpoint: "unit-test-hostname", Value: 0, CounterType: "GAUGE", Tags: "", Timestamp: time.Now().Unix(), Step: 60}},
	}
	for i, v := range tests {
		params := marshalStatsRow(v, nqmTarget, nqmAgent, 300, new(Fping))
		testTags := assembleTags(target, agent, tests[0])
		expecteds[i][3].Tags = testTags
		params[3].Tags = testTags
		//if !reflect.DeepEqual(expecteds[i][3], params[3]) {
		if expecteds[i][3].String() != params[3].String() {
			t.Error(expecteds[i][3], params[3])
		}
		t.Log(expecteds[i][3], params[3])
	}

	// tcpconn, 2 JSON parameters, only test the 2nd which is for Cassandra
	tests = []map[string]string{
		{"time": "13.6"},
		{"time": "-1"},
	}
	expecteds = [][]ParamToAgent{
		{{}, {Metric: "nqm-tcpconn", Endpoint: "unit-test-hostname", Value: 0, CounterType: "GAUGE", Tags: "", Timestamp: time.Now().Unix(), Step: 60}},
		{{}, {Metric: "nqm-tcpconn", Endpoint: "unit-test-hostname", Value: 0, CounterType: "GAUGE", Tags: "", Timestamp: time.Now().Unix(), Step: 60}},
	}
	for i, v := range tests {
		params := marshalStatsRow(v, nqmTarget, nqmAgent, 300, new(Tcpconn))
		testTags := assembleTags(target, agent, tests[0])
		expecteds[i][1].Tags = testTags
		params[1].Tags = testTags
		//if !reflect.DeepEqual(expecteds[i][1], params[1]) {
		if expecteds[i][1].String() != params[1].String() {
			t.Error(expecteds[i][1], params[1])
		}
		t.Log(expecteds[i][1], params[1])
	}
}
