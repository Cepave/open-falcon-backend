package check

import (
	"fmt"
	ot "github.com/Cepave/open-falcon-backend/common/reflect/types"
	"gopkg.in/check.v1"
	"reflect"
)

// Checks if obtained > expected
var LargerThan = &checkForCompare{largerThan, "LargerThan", "The left value **is not** larger than right value"}

// Checks if obtained >= expected
var LargerThanOrEqualTo = &checkForCompare{largerThanOrEqualTo, "LargerThanOrEqual", "The left value **is not** larger than or equal to right value"}

// Checks if obtained < expected
var SmallerThan = &checkForCompare{smallerThan, "SmallerThan", "The left value **is not** larger than right value"}

// Checks if obtained <= expected
var SmallerThanOrEqualTo = &checkForCompare{smallerThanOrEqualTo, "SmallerThanOrEqual", "The left value **is not** larger than or equal to right value"}

const (
	_ = iota
	largerThan
	largerThanOrEqualTo
	smallerThan
	smallerThanOrEqualTo
)

type checkForCompare struct {
	operator int
	name     string
	message  string
}

func (c *checkForCompare) Info() *check.CheckerInfo {
	return &check.CheckerInfo{
		Name:   c.name,
		Params: []string{"left", "right"},
	}
}
func (c *checkForCompare) Check(params []interface{}, names []string) (bool, string) {
	if performCompare(params[0], params[1], c.operator) {
		return true, ""
	}

	return false, c.message
}

func performCompare(left interface{}, right interface{}, operator int) bool {
	convertedLeftValue := convertNumberToCommonType(left)
	convertedRightValue := convertNumberToCommonType(right)

	switch leftValue := convertedLeftValue.(type) {
	case int64:
		return compareInt64(leftValue, convertedRightValue.(int64), operator)
	case uint64:
		return compareUint64(leftValue, convertedRightValue.(uint64), operator)
	case float64:
		return compareFloat64(leftValue, convertedRightValue.(float64), operator)
	}

	panic(fmt.Sprintf("Unsupported type for compare: %T", convertedLeftValue))
}

func convertNumberToCommonType(v interface{}) interface{} {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Convert(ot.TypeOfInt64).Interface()
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Convert(ot.TypeOfUint64).Interface()
	case float32, float64:
		return reflect.ValueOf(v).Convert(ot.TypeOfFloat64).Interface()
	}

	panic(fmt.Sprintf("Unsupported type: %T", v))
}

func compareInt64(left int64, right int64, operator int) bool {
	switch operator {
	case largerThan:
		return left > right
	case largerThanOrEqualTo:
		return left >= right
	case smallerThan:
		return left < right
	case smallerThanOrEqualTo:
		return left <= right
	}

	panic(fmt.Sprintf("Unknown operator: %d", operator))
}
func compareUint64(left uint64, right uint64, operator int) bool {
	switch operator {
	case largerThan:
		return left > right
	case largerThanOrEqualTo:
		return left >= right
	case smallerThan:
		return left < right
	case smallerThanOrEqualTo:
		return left <= right
	}

	panic(fmt.Sprintf("Unknown operator: %d", operator))
}
func compareFloat64(left float64, right float64, operator int) bool {
	switch operator {
	case largerThan:
		return left > right
	case largerThanOrEqualTo:
		return left >= right
	case smallerThan:
		return left < right
	case smallerThanOrEqualTo:
		return left <= right
	}

	panic(fmt.Sprintf("Unknown operator: %d", operator))
}
