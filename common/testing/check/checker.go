package check

import (
	"fmt"
	"gopkg.in/check.v1"
	"time"
)

var TimeEquals timeEqualsImpl

type timeEqualsImpl bool
var timeEqualsImplCheckInfo = &check.CheckerInfo {
	Name: "TimeEquals",
	Params: []string { "obtained", "expected" },
}

func (e timeEqualsImpl) Info() *check.CheckerInfo {
	return timeEqualsImplCheckInfo
}
func (e timeEqualsImpl) Check(params []interface{}, names []string) (result bool, errorMsg string) {
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
			if obtainedTimeValue.Unix() != params[i].(time.Time).Unix() {
				return false, "Checked values of time is not equal to obtained one"
			}
		default:
			return false, fmt.Sprintf("Type of obtained object(index:[%d]) is not time.Time: %v", i, t)
		}
	}

	return true, ""
}
