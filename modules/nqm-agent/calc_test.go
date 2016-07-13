package main

import (
	"reflect"
	"testing"
)

func TestConvToFloat(t *testing.T) {
	var nilSlice []float64
	tests := []struct {
		input    []string
		expected []float64
	}{
		{[]string{"13.24", "38.90", "19.62", "9.48", "13.62"}, []float64{13.24, 38.90, 19.62, 9.48, 13.62}},
		{[]string{"6.72", "29.08", "8.55", "7.40", "-", "6.26"}, []float64{6.72, 29.08, 8.55, 7.40, 6.26}},
		{[]string{"38.90", "1823.97", "10.48", "-", "25377.48"}, []float64{38.90, 10.48}},
		{[]string{"-", "-", "-"}, nilSlice},
	}
	for _, v := range tests {
		if !reflect.DeepEqual(convToFloat(v.input), v.expected) {
			t.Error(convToFloat(v.input), "!=", v.expected)
		}
		t.Log(convToFloat(v.input), v.expected)
	}
}

func TestCalcRow(t *testing.T) {
	// fping
	tests := [][]string{
		{"www.google.com", "13.24", "38.90", "19.62", "9.48", "13.62"},
		{"www.yahoo.com", "6.72", "29.08", "8.55", "7.40", "-", "6.26"},
		{"www.null.com", "-", "-", "-"},
	}

	// map[pkttransmit:5 pktreceive:5 rttmin:9.48 rttmax:38.90 rttavg:18.97 rttmdev:10.48 rttmedian:13.62]
	// map[rttmedian:7.40 pkttransmit:6 pktreceive:5 rttmin:6.26 rttmax:29.08 rttavg:11.60 rttmdev:8.77]

	expecteds := []map[string]string{
		{"rttmax": "38.90", "rttavg": "18.97", "rttmdev": "10.48", "rttmedian": "13.62", "pkttransmit": "5", "pktreceive": "5", "rttmin": "9.48"},
		{"rttmdev": "8.77", "rttmedian": "7.40", "pkttransmit": "6", "pktreceive": "5", "rttmin": "6.26", "rttmax": "29.08", "rttavg": "11.60"},
		{"rttmdev": "-1", "rttmedian": "-1", "pkttransmit": "3", "pktreceive": "0", "rttmin": "-1", "rttmax": "-1", "rttavg": "-1"},
	}
	fping := new(Fping)
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], calcRow(v, fping)) {
			t.Error(expecteds[i], calcRow(v, fping))
		}
		t.Log(calcRow(v, fping))
	}

	// tcpconn
	tests = [][]string{
		{"www.google.com", "13.24"},
		{"www.yahoo.com", "6.72"},
		{"www.null.com", "-"},
	}

	expecteds = []map[string]string{
		{"time": "13.24"},
		{"time": "6.72"},
		{"time": "-1"},
	}
	tcpconn := new(Tcpconn)
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], calcRow(v, tcpconn)) {
			t.Error(expecteds[i], calcRow(v, tcpconn))
		}
		t.Log(calcRow(v, tcpconn))
	}
}
