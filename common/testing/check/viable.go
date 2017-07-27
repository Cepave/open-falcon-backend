package check

import (
	"github.com/Cepave/open-falcon-backend/common/utils"
	"gopkg.in/check.v1"
	"reflect"
)

// See "utils.ValueExt.IsViable()" function
var ViableValue = viableValue(true)

type viableValue bool

func (v viableValue) Info() *check.CheckerInfo {
	return &check.CheckerInfo{
		Name:   "ViableValue",
		Params: []string{"obtained", "Need viable"},
	}
}
func (v viableValue) Check(params []interface{}, names []string) (bool, string) {
	needViable := params[1].(bool)

	checkedValue := reflect.ValueOf(params[0])
	valueExt := utils.ValueExt(checkedValue)

	valid := valueExt.IsViable()

	if needViable && !valid {
		return false, "Obtained value should not be nil"
	} else if !needViable && valid {
		return false, "Obtained value should be nil"
	}

	return true, ""
}
