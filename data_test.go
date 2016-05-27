package main

import (
	"reflect"
	"testing"
)

func TestTagsAssembler(t *testing.T) {
	agent := &nqmEndpointData{
		"-1", "-1", "-1", "-1", "-1",
	}
	target := &nqmEndpointData{
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
		if !reflect.DeepEqual(expecteds[i], TagsAssembler(target, agent, v)) {
			t.Error(expecteds[i], TagsAssembler(target, agent, v))
		}
		t.Log(TagsAssembler(target, agent, v))
	}
}

func init() {
	// Hostname is the config dependency which lies in func MarshalIntoParameters
	var cfg GeneralConfig
	generalConfig = &cfg
	cfg.Hostname = "unit-test-hostname"
}
