package main

import (
	"github.com/Cepave/common/model"
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
		{"rttmax": "46.5", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8.61"},
		{"rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "99", "rttmin": "5.42", "rttmax": "35.9", "rttavg": "10.9"},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], nqmParseFpingRow(v)) {
			t.Error(expecteds[i], nqmParseFpingRow(v))
		}
	}
}

func TestNqmFpingStat(t *testing.T) {
	tests := [][]string{
		{"www.google.com", "13.24", "38.90", "19.62", "9.48", "13.62"},
		{"www.yahoo.com", "6.72", "29.08", "8.55", "7.40", "-", "6.26"},
		{"www.null.com", "-", "-", "-"},
	}

	/*
	   data_test.go:37: map[pkttransmit:5 pktreceive:5 rttmin:9.48 rttmax:38.90 rttavg:18.97 rttmdev:10.48 rttmedian:13.62]
	   data_test.go:37: map[rttmedian:7.40 pkttransmit:6 pktreceive:5 rttmin:6.26 rttmax:29.08 rttavg:11.60 rttmdev:8.77]
	*/
	expecteds := []map[string]string{
		{"rttmax": "38.90", "rttavg": "18.97", "rttmdev": "10.48", "rttmedian": "13.62", "pkttransmit": "5", "pktreceive": "5", "rttmin": "9.48"},
		{"rttmdev": "8.77", "rttmedian": "7.40", "pkttransmit": "6", "pktreceive": "5", "rttmin": "6.26", "rttmax": "29.08", "rttavg": "11.60"},
		{"rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "3", "pktreceive": "0", "rttmin": "-1", "rttmax": "-1", "rttavg": "-1"},
	}

	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], nqmFpingStat(v)) {
			t.Error(expecteds[i], nqmFpingStat(v))
		}
		t.Log(nqmFpingStat(v))
	}
}

func TestParseFpingRow(t *testing.T) {
	tests := []string{
		"www.google.com : xmt/rcv/%loss = 100/100/0%, min/avg/max = 8.61/14.5/46.5",
		"www.yahoo.com  : xmt/rcv/%loss = 100/99/1%, min/avg/max = 5.42/10.9/35.9",
		"www.google.com : 13.24 38.90 19.62 9.48 13.62",
		"www.yahoo.com : 6.72 29.08 8.55 7.40 - 6.26",
	}

	expecteds := [][]string{
		{"www.google.com", "xmt", "rcv", "loss", "100", "100", "0", "min", "avg", "max", "8.61", "14.5", "46.5"},
		{"www.yahoo.com", "xmt", "rcv", "loss", "100", "99", "1", "min", "avg", "max", "5.42", "10.9", "35.9"},
		{"www.google.com", "13.24", "38.90", "19.62", "9.48", "13.62"},
		{"www.yahoo.com", "6.72", "29.08", "8.55", "7.40", "-", "6.26"},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], parseFpingRow(v)) {
			t.Error(expecteds[i], parseFpingRow(v))
		}
		t.Log(parseFpingRow(v))
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
		{"rttmax": "46.5", "rttavg": "14.5", "rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "100", "pktreceive": "100", "rttmin": "8.61"},
	}

	expecteds := []string{
		"agent-id=-1,agent-isp-id=-1,agent-province-id=-1,agent-city-id=-1,agent-name-tag-id=-1,target-id=-2,target-isp-id=-2,target-province-id=-2,target-city-id=-2,target-name-tag-id=-2,rttmin=8.61,rttmax=46.5,rttavg=14.5,rttmdev=-1,rttmedian=-1,pkttransmit=100,pktreceive=100",
	}

	t_out := nqmTagsAssembler(target, agent, tests[0])
	if t_out != expecteds[0] {
		t.Error(expecteds[0], t_out)
	}
}

func init() {
	/*var (
	  agentData *nqmEndpointData
	  )
	*/
	// agentData is initialized for testMarshalIntoParameters.
	agentData = &nqmEndpointData{
		Id:         "1000",
		IspId:      "1001",
		ProvinceId: "1002",
		CityId:     "1003",
		NameTagId:  "1004",
	}
}

func TestMarshalIntoParameters(t *testing.T) {
	/*
	   www.google.com : 21.77 15.88 14.57 18.94 17.23
	   www.yahoo.com  : 12.86 7.67 6.81 6.96 8.65
	*/
	tests := []string{
		"www.google.com : 13.24 38.90 19.62 9.48 13.62",
		"www.yahoo.com : 6.72 29.08 8.55 7.40 - 6.26",
	}
	var test_target_list []model.NqmTarget
	for i, _ := range tests {
		t := model.NqmTarget{
			Id:           i,
			Host:         "test host",
			IspId:        int16(i),
			IspName:      "test isp",
			ProvinceId:   int16(i),
			ProvinceName: "test province",
			CityId:       int16(i),
			CityName:     "test city",
			NameTag:      "test nametag",
		}
		test_target_list = append(test_target_list, t)
	}

	out := MarshalIntoParameters(tests, test_target_list)
	t.Log(out)
	//t.Log(test_target_list)
}
