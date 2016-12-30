package utils

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

const (
	DefaultDirection byte = 0
	// Sorting by ascending
	Ascending byte = 1
	// Sorting by descending
	Descending byte = 2

	// The sequence should be higher
	SeqHigher = 1
	// The sequence is equal
	SeqEqual = 0
	// The sequence should be lower
	SeqLower = -1
)

// Compares any value. supporting:
//
// String
// Interger
// Unsigned Integer
// Float
//
// net.IP
//
// This funcation would be slower than typed function
func CompareAny(left interface{}, right interface{}, direction byte) int {
	valueOfLeft, valueOfRight := reflect.ValueOf(left), reflect.ValueOf(right)
	if valueOfLeft.Type() != valueOfRight.Type() {
		panic(fmt.Sprintf("Types for both side are not same. Left: %T. Right: %T", left, right))
	}

	switch leftValue := left.(type) {
	case string:
		return CompareString(leftValue, right.(string), direction)
	case int:
		return CompareInt(int64(leftValue), int64(right.(int)), direction)
	case int8:
		return CompareInt(int64(leftValue), int64(right.(int8)), direction)
	case int16:
		return CompareInt(int64(leftValue), int64(right.(int16)), direction)
	case int32:
		return CompareInt(int64(leftValue), int64(right.(int32)), direction)
	case int64:
		return CompareInt(leftValue, right.(int64), direction)
	case uint:
		return CompareUint(uint64(leftValue), uint64(right.(uint)), direction)
	case uint8:
		return CompareUint(uint64(leftValue), uint64(right.(uint8)), direction)
	case uint16:
		return CompareUint(uint64(leftValue), uint64(right.(uint16)), direction)
	case uint32:
		return CompareUint(uint64(leftValue), uint64(right.(uint32)), direction)
	case uint64:
		return CompareUint(leftValue, right.(uint64), direction)
	case float32:
		return CompareFloat(float64(leftValue), float64(right.(float32)), direction)
	case float64:
		return CompareFloat(leftValue, right.(float64), direction)
	case net.IP:
		return CompareIpAddress(leftValue, right.(net.IP), direction)
	}

	panic(fmt.Sprintf("Unsupported type for comparison on left value: %T", left))
}

// Compares IP address from left-most byte(numeric-sensitive)
func CompareIpAddress(leftIp net.IP, rightIp net.IP, direction byte) int {
	if r, hasNil := CompareNil(leftIp, rightIp, direction); hasNil {
		return r
	}

	leftIp16, rightIp16 := leftIp.To16(), rightIp.To16()

	for i := 0; i < 16; i++ {
		result := CompareUint(uint64(leftIp16[i]), uint64(rightIp16[i]), direction)

		if result != SeqEqual {
			return result
		}
	}

	return SeqEqual
}

// For ascending:
//
// Nil / Not Nil: SeqHigher
// Nil / Nil: SeqEqual
// Not Nil / Nil: SeqLower
//
// descending is the Reverse of ascending
//
// Returns:
// 	if the second value is true, means at least one of their value is nil
func CompareNil(left interface{}, right interface{}, direction byte) (r int, hasNil bool) {
	valueOfLeft, valueOfRight := reflect.ValueOf(left), reflect.ValueOf(right)

	hasNil = valueOfLeft.IsNil() || valueOfRight.IsNil()

	r = SeqEqual
	if valueOfLeft.IsNil() && !valueOfRight.IsNil() {
		r = SeqHigher
	} else if !valueOfLeft.IsNil() && valueOfRight.IsNil() {
		r = SeqLower
	}

	r = ReverseIfDescending(r, direction)
	return
}

func CompareString(left string, right string, direction byte) int {
	r := -strings.Compare(left, right)
	return ReverseIfDescending(r, direction)
}

func CompareInt(left int64, right int64, direction byte) int {
	if left == right {
		return SeqEqual
	}

	r := SeqLower
	if left < right {
		r = SeqHigher
	}

	return ReverseIfDescending(r, direction)
}
func CompareUint(left uint64, right uint64, direction byte) int {
	if left == right {
		return SeqEqual
	}

	r := SeqLower
	if left < right {
		r = SeqHigher
	}

	return ReverseIfDescending(r, direction)
}
func CompareFloat(left float64, right float64, direction byte) int {
	if left == right {
		return SeqEqual
	}

	r := SeqLower
	if left < right {
		r = SeqHigher
	}

	return ReverseIfDescending(r, direction)
}

func ReverseIfDescending(result int, direction byte) int {
	if direction == Descending {
		return -result
	}

	return result
}
