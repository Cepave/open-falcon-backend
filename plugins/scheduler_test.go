package plugins

import (
	"reflect"
	"testing"
)

var in []string

func init() {
	in = []string{
		"./test/scheduler_test_no_exec.py",
		"./test/scheduler_test_no_shebang.py",
		"./test/scheduler_test.py",
	}
}

func TestNoOwnerExecPerm(t *testing.T) {
	expect := []bool{true, true, false}
	for i, v := range in {
		real := noOwnerExecPerm(v)
		if expect[i] != real {
			t.Error("Input value is:", v, "Expected value is:", expect[i], ", Real value is:", real)
		}
	}

}

func TestHasShebang(t *testing.T) {
	expect := []bool{true, false, true}
	for i, v := range in {
		real := hasShebang(v)
		if expect[i] != real {
			t.Error("Input value is:", v, "Expected value is:", expect[i], ", Real value is:", real)
		}
	}
}

func TestGetInterpreterCmd(t *testing.T) {
	expect := []string{
		"/usr/bin/python",
		"-O",
		"./test/scheduler_test.py",
	}

	real := getInterpreterCmd(in[2])
	if !reflect.DeepEqual(expect, real) {
		t.Error("Input value is:", in[2], "Expected value is:", expect, "length of value is", len(expect))
		t.Error("Real value is:", real, "length of real value is:", len(real))
	}
}
