package utils

import (
	"reflect"
	"sort"
)

// Checks the two array of strings are same
//
// This function would sort there two strings first, then check their equality.
//
// The array of nil and array of empty are seen the same one.
func AreArrayOfStringsSame(leftArray []string, rightArray []string) bool {
	leftArrayOfStrings := make([]string, len(leftArray))
	rightArrayOfStrings := make([]string, len(rightArray))

	copy(leftArrayOfStrings, leftArray)
	copy(rightArrayOfStrings, rightArray)

	sort.Strings(leftArrayOfStrings)
	sort.Strings(rightArrayOfStrings)

	return reflect.DeepEqual(leftArrayOfStrings, rightArrayOfStrings)
}
