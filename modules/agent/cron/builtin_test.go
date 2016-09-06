package cron

import (
	"reflect"
	"testing"

	"github.com/Cepave/open-falcon-backend/common/model"
)

// test line number 125~130
func TestOWL700(t *testing.T) {
	// initializing
	var metric model.BuiltinMetric
	metric.Tags = "proc.num"
	var procs = make(map[string]map[int]string)
	var tmpMap = make(map[int]string)
	var emptyMap = make(map[int]string)

	// begin to test
	// test no any tags
	t.Log("test no tags")
	tmpMap = make(map[int]string)
	tmpMap[4] = "placeholder"
	procs[metric.Tags] = emptyMap
	if true {
		_, nameExist := tmpMap[1]
		_, cmdExist := tmpMap[2]
		_, wrongTagsExist := tmpMap[3]
		if !wrongTagsExist && !(nameExist && cmdExist) {
			procs[metric.Tags] = tmpMap
		}
	}
	if !reflect.DeepEqual(tmpMap, procs[metric.Tags]) {
		t.Error("must be equal")
	}
	t.Log("inside tmpMap is:", tmpMap)
	t.Log("inside procs[metric.Tags] is:", procs[metric.Tags])
	// test name case
	t.Log("test name case")
	tmpMap = make(map[int]string)
	tmpMap[4] = "placeholder"
	procs[metric.Tags] = emptyMap
	tmpMap[1] = "proc.num/name"
	if true {
		_, nameExist := tmpMap[1]
		_, cmdExist := tmpMap[2]
		_, wrongTagsExist := tmpMap[3]
		if !wrongTagsExist && !(nameExist && cmdExist) {
			procs[metric.Tags] = tmpMap
		}
	}
	if !reflect.DeepEqual(tmpMap, procs[metric.Tags]) {
		t.Error("must be equal")
	}
	t.Log("inside tmpMap is:", tmpMap)
	t.Log("inside procs[metric.Tags] is:", procs[metric.Tags])
	// test cmd case
	t.Log("test cmd case")
	tmpMap = make(map[int]string)
	tmpMap[4] = "placeholder"
	procs[metric.Tags] = emptyMap
	tmpMap[2] = "proc.num/cmdline"
	if true {
		_, nameExist := tmpMap[1]
		_, cmdExist := tmpMap[2]
		_, wrongTagsExist := tmpMap[3]
		if !wrongTagsExist && !(nameExist && cmdExist) {
			procs[metric.Tags] = tmpMap
		}
	}
	if !reflect.DeepEqual(tmpMap, procs[metric.Tags]) {
		t.Error("must be equal")
	}
	t.Log("inside tmpMap is:", tmpMap)
	t.Log("inside procs[metric.Tags] is:", procs[metric.Tags])
	// test cmd and name both exist case
	t.Log("test cmd and name both exist case")
	tmpMap = make(map[int]string)
	tmpMap[4] = "placeholder"
	procs[metric.Tags] = emptyMap
	tmpMap[1] = "proc.num/name"
	tmpMap[2] = "proc.num/cmdline"
	if true {
		_, nameExist := tmpMap[1]
		_, cmdExist := tmpMap[2]
		_, wrongTagsExist := tmpMap[3]
		if !wrongTagsExist && !(nameExist && cmdExist) {
			procs[metric.Tags] = tmpMap
		}
	}
	if reflect.DeepEqual(tmpMap, procs[metric.Tags]) {
		t.Error("must be not equal")
	}
	t.Log("inside tmpMap is:", tmpMap)
	t.Log("inside procs[metric.Tags] is:", procs[metric.Tags])
	// test wrong case
	t.Log("test wrong case")
	tmpMap = make(map[int]string)
	tmpMap[4] = "placeholder"
	procs[metric.Tags] = emptyMap
	tmpMap[3] = "proc.num/number"
	if true {
		_, nameExist := tmpMap[1]
		_, cmdExist := tmpMap[2]
		_, wrongTagsExist := tmpMap[3]
		if !wrongTagsExist && !(nameExist && cmdExist) {
			procs[metric.Tags] = tmpMap
		}
	}
	if reflect.DeepEqual(tmpMap, procs[metric.Tags]) {
		t.Error("must be not equal")
	}
	t.Log("inside tmpMap is:", tmpMap)
	t.Log("inside procs[metric.Tags] is:", procs[metric.Tags])
}
