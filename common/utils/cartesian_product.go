package utils

import (
	"reflect"

	"github.com/juju/errors"
)

// Performs Cartesian product and gives list of product results
//
// Every element should be slice or array.
func CartesianProduct(listOfSets ...interface{}) [][]interface{} {
	numberOfSets := len(listOfSets)

	if numberOfSets == 0 {
		return ([][]interface{})(nil)
	}

	currentSet := reflect.ValueOf(listOfSets[0])

	/**
	 * Checks type of set(must be slice or array)
	 */
	currentKind := currentSet.Kind()
	if currentKind != reflect.Array &&
		currentKind != reflect.Slice {
		panic(errors.Details(
			errors.New("Element must be slice or array"),
		))
	}
	// :~)

	sizeOfCurrentSet := currentSet.Len()

	result := make([][]interface{}, 0)

	for i := 0; i < sizeOfCurrentSet; i++ {
		currentElement := currentSet.Index(i).Interface()

		if numberOfSets == 1 {
			result = append(result, []interface{}{currentElement})
			continue
		}

		/**
		 * Appends product of rest sets(recursively)
		 */
		subProduct := CartesianProduct(listOfSets[1:numberOfSets]...)
		sizeOfSubProduct := len(subProduct)

		for j := 0; j < sizeOfSubProduct; j++ {
			currentRecord := []interface{}{currentElement}
			currentRecord = append(currentRecord, subProduct[j]...)

			result = append(result, currentRecord)
		}
		// :~)
	}

	return result
}
