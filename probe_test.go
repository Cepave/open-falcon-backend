package main

import (
	"testing"
)

// fping -p 20 -i 10 -C 5 -a www.google.com www.yahoo.com
func TestProbe(t *testing.T) {
	tests := []string{
		"fping", "-p", "20", "-i", "10", "-C", "5", "-a",
		"www.google.com",
		"www.yahoo.com",
	}

	t.Log("input cmd is:", tests)
	expected := Probe(tests)
	for _, v := range expected {
		t.Log(v)
	}
}
