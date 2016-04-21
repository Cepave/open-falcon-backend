package main

import (
	"reflect"
	"testing"
)

/*
www.google.com : xmt/rcv/%loss = 100/100/0%, min/avg/max = 8.61/14.5/46.5
www.yahoo.com  : xmt/rcv/%loss = 100/99/1%, min/avg/max = 5.42/10.9/35.9
*/

func TestNqmParseFpingRow(t *testing.T) {
	tests := [][]string{
		{"www.google.com", "xmt", "rcv", "loss", "100", "100", "0", "min", "avg", "max", "8.61", "14.5", "46.5"},
		{"www.yahoo.com", "xmt", "rcv", "loss", "100", "99", "1", "min", "avg", "max", "5.42", "10.9", "35.9"},
	}

	expecteds := []map[string]string{
		{"rttmax": "46", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8"},
		{"rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "99", "rttmin": "5", "rttmax": "35", "rttavg": "10.9"},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], nqmParseFpingRow(v)) {
			t.Error(expecteds[i], nqmParseFpingRow(v))
		}
	}
}

func TestParseFpingRow(t *testing.T) {
	tests := []string{
		"www.google.com : xmt/rcv/%loss = 100/100/0%, min/avg/max = 8.61/14.5/46.5",
		"www.yahoo.com  : xmt/rcv/%loss = 100/99/1%, min/avg/max = 5.42/10.9/35.9",
	}

	expecteds := [][]string{
		{"www.google.com", "xmt", "rcv", "loss", "100", "100", "0", "min", "avg", "max", "8.61", "14.5", "46.5"},
		{"www.yahoo.com", "xmt", "rcv", "loss", "100", "99", "1", "min", "avg", "max", "5.42", "10.9", "35.9"},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], parseFpingRow(v)) {
			t.Error(expecteds[i], parseFpingRow(v))
		}
	}
}

func TestNqmTagsAssembler(t *testing.T) {
	agent := &nqmEndpointData{
		"-1", "-1", "-1", "-1", "-1",
	}
	target := &nqmEndpointData{
		"-2", "-2", "-2", "-2", "-2",
	}
	tests := []map[string]string{
		{"rttmax": "46", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8"},
	}

	expecteds := []string{
		"agent-id=-1,agent-isp-id=-1,agent-province-id=-1,agent-city-id=-1,agent-name-tag-id=-1,target-id=-2,target-isp-id=-2,target-province-id=-2,target-city-id=-2,target-name-tag-id=-2,rttmin=8,rttmax=46,rttavg=14.5,rttmdev=-1,rttmedian=-1,pkttransmit=100,pktreceive=100",
	}

	t_out := nqmTagsAssembler(target, agent, tests[0])
	if t_out != expecteds[0] {
		t.Error(expecteds[0], t_out)
	}
}
