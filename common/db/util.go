package db

import (
	"database/sql"
	"fmt"
	csc "github.com/Cepave/open-falcon-backend/common/strconv"
	"regexp"
	"strconv"
	"strings"
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
		currentValue, err := strconv.ParseUint(v, 10, 8)
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
// separator - The separator
func GroupedStringToIntArray(groupedValue sql.NullString, separator string) []int64 {
	if !groupedValue.Valid {
		return make([]int64, 0)
	}

	return csc.SplitStringToIntArray(groupedValue.String, separator)
}

// Converts grouped string to array of uint64
//
// groupedValue - The grouped value retrieved from database, could be null
// separator - The separator
func GroupedStringToUintArray(groupedValue sql.NullString, separator string) []uint64 {
	if !groupedValue.Valid {
		return make([]uint64, 0)
	}

	return csc.SplitStringToUintArray(groupedValue.String, separator)
}

// Converts grouped string to array of string
//
// groupedValue - The grouped value retrieved from database, could be null
// separator - The separator
func GroupedStringToStringArray(groupedValue sql.NullString, separator string) []string {
	if !groupedValue.Valid {
		return make([]string, 0)
	}

	return strings.Split(groupedValue.String, separator)
}
