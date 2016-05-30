package main

import (
	"strings"
	"testing"
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
