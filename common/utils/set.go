package utils

import (
	"reflect"
)

var valueOfTrue = reflect.ValueOf(true)
var typeOfBoolean = valueOfTrue.Type()

// Uses reflect.* to make unique element of input array
func UniqueElements(arrayOrSlice interface{}) interface{} {
	if arrayOrSlice == nil {
		return nil
	}

	uniqueArrayType := reflect.TypeOf(arrayOrSlice)

	arrayObject := MakeAbstractArray(arrayOrSlice).
		FilterWith(NewUniqueFilter(uniqueArrayType))

	return arrayObject.GetArray()
}

// Makes the array of strings unique, which is stable(the sequence of output is same as input)
func UniqueArrayOfStrings(arrayOfStrings []string) []string {
	if arrayOfStrings == nil {
		return nil
	}

	result := UniqueElements(arrayOfStrings)
	return result.([]string)
}
