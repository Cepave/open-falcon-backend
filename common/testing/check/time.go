package check

import (
	"fmt"
	"gopkg.in/check.v1"
	"time"
)

var TimeEquals timeEqualsImpl
var TimeBefore timeBeforeImpl
var TimeAfter timeAfterImpl

type timeEqualsImpl bool
func (e timeEqualsImpl) Info() *check.CheckerInfo {
	return &check.CheckerInfo {
		Name: "TimeEquals",
		Params: []string { "obtained", "expected" },
	}
}
func (e timeEqualsImpl) Check(params []interface{}, names []string) (result bool, errorMsg string) {
	return checkForTimes(
		params, names,
		func(firstValue time.Time, secondValue time.Time) (bool, string) {
			if firstValue.Unix() != secondValue.Unix() {
				return false, "Objtained time is not equal to expected one"
			}

			return true, ""
		},
	)
}

type timeBeforeImpl bool
func (e timeBeforeImpl) Info() *check.CheckerInfo {
	return &check.CheckerInfo {
		Name: "TimeBefore",
		Params: []string { "obtained", "newer time" },
	}
}
func (e timeBeforeImpl) Check(params []interface{}, names []string) (result bool, errorMsg string) {
	return checkForTimes(
		params, names,
		func(firstValue time.Time, secondValue time.Time) (bool, string) {
			if !firstValue.Before(secondValue) {
				return false, "Obtained time is not before the newer one"
			}

			return true, ""
		},
	)
}

type timeAfterImpl bool
func (e timeAfterImpl) Info() *check.CheckerInfo {
	return &check.CheckerInfo {
		Name: "TimeAfter",
		Params: []string { "obtained", "older time" },
	}
}
func (e timeAfterImpl) Check(params []interface{}, names []string) (result bool, errorMsg string) {
	return checkForTimes(
		params, names,
		func(firstValue time.Time, secondValue time.Time) (bool, string) {
			if !firstValue.After(secondValue) {
				return false, "Obtained time is not after older one"
			}

			return true, ""
		},
	)
}

type checkTimeFunc func(firstValue time.Time, secondValue time.Time) (bool, string)

func checkForTimes(
	params []interface{}, names []string,
	checkImpl checkTimeFunc,
) (result bool, errorMsg string) {
	var obtainedTimeValue time.Time

	switch t := params[0].(type) {
	case time.Time:
		obtainedTimeValue = params[0].(time.Time)
	default:
		return false, fmt.Sprintf("Type of obtained object is not time.Time: %v", t)
	}

	for i := 1; i < len(params); i++ {
		switch t := params[i].(type) {
		case time.Time:
			return checkImpl(obtainedTimeValue, params[i].(time.Time))
		default:
			return false, fmt.Sprintf("Type of obtained object(index:[%d]) is not time.Time: %v", i, t)
		}
	}

	return true, ""
}
