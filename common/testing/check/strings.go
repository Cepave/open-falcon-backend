package check

import (
	"gopkg.in/check.v1"
	"strings"
)

var StringContains = stringContains(true)

type stringContains bool

func (s stringContains) Info() *check.CheckerInfo {
	return &check.CheckerInfo{
		Name:   "String contains",
		Params: []string{"obtained", "checked"},
	}
}
func (c stringContains) Check(params []interface{}, names []string) (bool, string) {
	checkedValue := params[0].(string)
	containedValue := params[1].(string)

	if strings.Contains(checkedValue, containedValue) {
		return true, ""
	}

	return false, "Obtained value doesn't contain the checked value"
}
