package core

import (
	"testing"
)

func TestToInt64(t *testing.T) {
	_, e := ToInt64(1)
	if e != nil {
		t.Errorf("ToInt64:\n Expect => %v\n Got    => %v\n", nil, e)
	}

	_, e = ToInt64("1")
	if e == nil {
		t.Errorf("ToInt64:\n Expect => %v\n Got    => %v\n", "error", e)
	}
}
