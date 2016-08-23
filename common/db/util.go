package db

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
)

// Converts grouped string to array of int64
//
// groupedValue - The grouped value retrieved from database, could be null
// seperator - The seperator
func GroupedStringToIntArray(groupedValue sql.NullString, seperator string) []int64 {
	stringArray := GroupedStringToStringArray(groupedValue, seperator)
	if stringArray == nil {
		return nil
	}
	result := make([]int64, len(stringArray))

	for i, v := range stringArray {
		var intValue int64
		var parseError error

		if intValue, parseError = strconv.ParseInt(v, 10, 64)
			parseError != nil {
			panic(fmt.Errorf("Cannot parse value in array to Int. Index: [%d]. Value: [%v]", i, v))
		}

		result[i] = intValue
	}

	return result
}

// Converts grouped string to array of uint64
//
// groupedValue - The grouped value retrieved from database, could be null
// seperator - The seperator
func GroupedStringToUintArray(groupedValue sql.NullString, seperator string) []uint64 {
	stringArray := GroupedStringToStringArray(groupedValue, seperator)
	if stringArray == nil {
		return nil
	}

	result := make([]uint64, len(stringArray))

	for i, v := range stringArray {
		var uintValue uint64
		var parseError error

		if uintValue, parseError = strconv.ParseUint(v, 10, 64)
			parseError != nil {
			panic(fmt.Errorf("Cannot parse value in array to Uint. Index: [%d]. Value: [%v]", i, v))
		}

		result[i] = uintValue
	}

	return result
}

// Converts grouped string to array of string
//
// groupedValue - The grouped value retrieved from database, could be null
// seperator - The seperator
func GroupedStringToStringArray(groupedValue sql.NullString, seperator string) []string {
	if !groupedValue.Valid {
		return nil
	}

	return strings.Split(groupedValue.String, seperator)
}
