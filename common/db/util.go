package db

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"strconv"
)

var ipV4RegExp = regexp.MustCompile("^(?:[0-9]{1,3})(?:\\.[0-9]{1,3})?(?:\\.[0-9]{1,3})?(?:\\.[0-9]{1,3})?$")

// Checks whether or not the string is partial IPv4:
//
// 	123
// 	123.74
// 	123.74.109.81
func IsPartialIpV4(ip string) bool {
	return ipV4RegExp.MatchString(ip)
}

// Converts IPv4 address of string to byte array could be used in SQL.
func IpV4ToBytesForLike(ip string) ([]byte, error) {
	if !IsPartialIpV4(ip) {
		return nil, fmt.Errorf("IP [%s] doesn't match IPv4 address")
	}

	splitedIps := strings.Split(ip, ".")

	var finalIp []byte
	for _, v := range splitedIps {
		if v == "" {
			continue
		}

		/**
		 * Converts value to byte
		 */
		currentValue, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse [%s] for IPv4: %v", v, err)
		}
		// :~)

		/**
		 * Puts '\' to escape '%' in binary
		 */
		if currentValue == 0x25 {
			finalIp = append(finalIp, 0x5C)
		}
		// :~)

		finalIp = append(finalIp, byte(currentValue))
	}

	finalIp = append(finalIp, 0x25)

	return finalIp, nil
}

// Converts grouped string to array of int64
//
// groupedValue - The grouped value retrieved from database, could be null
// seperator - The seperator
func GroupedStringToIntArray(groupedValue sql.NullString, seperator string) []int64 {
	if !groupedValue.Valid {
		return nil
	}

	return GroupedPlainStringToIntArray(groupedValue.String, seperator)
}

// Converts grouped string to array of int64
//
// groupedValue - The grouped value retrieved from database, could be null
// seperator - The seperator
func GroupedPlainStringToIntArray(groupedValue string, seperator string) []int64 {
	if groupedValue == "" {
		return nil
	}

	stringArray := strings.Split(groupedValue, seperator)
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
	if !groupedValue.Valid {
		return nil
	}

	return GroupedPlainStringToUintArray(groupedValue.String, seperator)
}
func GroupedPlainStringToUintArray(groupedValue string, seperator string) []uint64 {
	if groupedValue == "" {
		return nil
	}

	stringArray := strings.Split(groupedValue, seperator)
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
