package funcs

import "testing"

func TestContainsCollector(t *testing.T) {
	// cat /proc/net/dev
	in := []string{"eth0", "eth1", "docker0", "em0", "bond3"}
	ifacePrefix := []string{"eth", "em", "bond", "enp"}
	out := []bool{true, true, false, true, true}
	for i, val := range in {
		if out[i] != containsCollector(val, ifacePrefix) {
			t.Error("unepected result: ", out[i], val)
		}
	}
}
