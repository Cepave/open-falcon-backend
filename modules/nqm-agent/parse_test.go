package main

import (
	"reflect"
	"testing"
)

func TestParseRow(t *testing.T) {
	tests := []string{
		"www.google.com   : 7.23 6.93 5.73 5.60",
		"www.youtube.com  : 4.55 6.60 6.61 7.99",
		"www.facebook.com : 64.82 69.21 70.67 70.90",
		"210.242.127.93  : 7.08 9.59 7.51 7.16",
		"210.242.127.118 : 8.46 5.56 8.39 6.46",
		"31.13.95.36     : 95.41 97.45 100.51 102.43",
	}

	expecteds := [][]string{
		{"www.google.com", "7.23", "6.93", "5.73", "5.60"}, {"www.youtube.com", "4.55", "6.60", "6.61", "7.99"}, {"www.facebook.com", "64.82", "69.21", "70.67", "70.90"},
		{"210.242.127.93", "7.08", "9.59", "7.51", "7.16"}, {"210.242.127.118", "8.46", "5.56", "8.39", "6.46"}, {"31.13.95.36", "95.41", "97.45", "100.51", "102.43"},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], parseRow(v)) {
			t.Error(expecteds[i], parseRow(v))
		}
		t.Log(expecteds[i], parseRow(v))
	}
}
func TestParse(t *testing.T) {
	tests := [][]string{
		{"www.google.com   : 7.23 6.93 5.73 5.60", "www.youtube.com  : 4.55 6.60 6.61 7.99", "www.facebook.com : 64.82 69.21 70.67 70.90"},
		{"210.242.127.93  : 7.08 9.59 7.51 7.16", "210.242.127.118 : 8.46 5.56 8.39 6.46", "31.13.95.36     : 95.41 97.45 100.51 102.43"},
	}

	expecteds := [][][]string{
		{{"www.google.com", "7.23", "6.93", "5.73", "5.60"}, {"www.youtube.com", "4.55", "6.60", "6.61", "7.99"}, {"www.facebook.com", "64.82", "69.21", "70.67", "70.90"}},
		{{"210.242.127.93", "7.08", "9.59", "7.51", "7.16"}, {"210.242.127.118", "8.46", "5.56", "8.39", "6.46"}, {"31.13.95.36", "95.41", "97.45", "100.51", "102.43"}},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(expecteds[i], Parse(v)) {
			t.Error(expecteds[i], Parse(v))
		}
		t.Log(expecteds[i], Parse(v))
	}
}
