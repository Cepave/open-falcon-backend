package strconv

import (
	"fmt"
	"strconv"
	"strings"
)

func SplitStringToIntArray(values string, separator string) []int64 {
	if values == "" {
		return make([]int64, 0)
	}

	stringArray := strings.Split(values, separator)
	result := make([]int64, len(stringArray))

	for i, v := range stringArray {
		var intValue int64
		var parseError error

		if intValue, parseError = strconv.ParseInt(v, 10, 64); parseError != nil {
			panic(fmt.Errorf("Cannot parse value in array to Int. Index: [%d]. Value: [%v]", i, v))
		}

		result[i] = intValue
	}

	return result
}
func SplitStringToUintArray(values string, separator string) []uint64 {
	if values == "" {
		return make([]uint64, 0)
	}

	stringArray := strings.Split(values, separator)
	result := make([]uint64, len(stringArray))

	for i, v := range stringArray {
		var uintValue uint64
		var parseError error

		if uintValue, parseError = strconv.ParseUint(v, 10, 64); parseError != nil {
			panic(fmt.Errorf("Cannot parse value in array to Uint. Index: [%d]. Value: [%v]", i, v))
		}

		result[i] = uintValue
	}

	return result
}
