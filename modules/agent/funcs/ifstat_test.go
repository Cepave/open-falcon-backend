package funcs

import "testing"

func TestContainsCollector(t *testing.T) {
	// cat /proc/net/dev
	in := []string{"eth0", "eth1", "docker0", "em0", "bond3"}
	ethAll := []string{"eth", "em", "bond", "enp"}
	out := []bool{true, true, false, true, true}
	for i, val := range in {
		if out[i] != containsCollector(val, ethAll) {
			t.Error("unepected result: ", out[i], val)
		}
	}
}

func TestCoreNetMetrics(t *testing.T) {
	ifacePrefix := []string{"eth", "lo", "bond", "em", "br"}
	ethAll := []string{"eth", "lo"}
	t.Log("Test eth_all: ", ethAll)
	for _, val := range CoreNetMetrics(ifacePrefix, ethAll) {
		if val.Metric == "net.if.in.bits" || val.Metric == "net.if.out.bits" {
			t.Log(val)
		}
	}
}
