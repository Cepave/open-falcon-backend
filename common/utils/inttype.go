package utils

import (
	"reflect"
	"sort"
)

func integerToAny(source interface{}, typeOfTarget reflect.Type) interface{} {
	if source == nil {
		return nil
	}

	arrayObject := MakeAbstractArray(source).
		MapTo(IdentityMapper, typeOfTarget)

	return arrayObject.GetArray()
}

// Converts 64 bits integer(unsigned) to 32 bits one
func UintTo32(source []uint64) []uint32 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfUint32)
	return convertedArray.([]uint32)
}

// Converts 64 bits integer(unsigned) to 16 bits one
func UintTo16(source []uint64) []uint16 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfUint16)
	return convertedArray.([]uint16)
}

// Converts 64 bits integer(unsigned) to 8 bits one
func UintTo8(source []uint64) []uint8 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfUint8)
	return convertedArray.([]uint8)
}

// Converts 64 bits integer to 32 bits one
func IntTo32(source []int64) []int32 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfInt32)
	return convertedArray.([]int32)
}

// Converts 64 bits integer to 16 bits one
func IntTo16(source []int64) []int16 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfInt16)
	return convertedArray.([]int16)
}

// Converts 64 bits integer to 8 bits one
func IntTo8(source []int64) []int8 {
	if source == nil {
		return nil
	}

	convertedArray := integerToAny(source, TypeOfInt8)
	return convertedArray.([]int8)
}

type Int64Slice []int64

func (s Int64Slice) Len() int           { return len(s) }
func (s Int64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Int64Slice) Less(i, j int) bool { return s[i] < s[j] }

type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }

// Sorts the value of int64 with unique processing
func SortAndUniqueInt64(source []int64) []int64 {
	if source == nil {
		return nil
	}

	uniqueElements := UniqueElements(source).([]int64)
	sort.Sort(Int64Slice(uniqueElements))

	return uniqueElements
}

// Sorts the value of int64 with unique processing
func SortAndUniqueUint64(source []uint64) []uint64 {
	if source == nil {
		return nil
	}

	uniqueElements := UniqueElements(source).([]uint64)
	sort.Sort(Uint64Slice(uniqueElements))

	return uniqueElements
}
