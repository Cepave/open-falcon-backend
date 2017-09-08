package validate

import (
	rf "reflect"

	v "gopkg.in/go-playground/validator.v9"
)

const (
	// Supports slice or array
	//
	// Validatation would fail if all of elements in slice(array) are zero.
	//
	// Nil or empty slice is passed.
	TagNonZeroSlice = "non_zero_slice"
)

func NonZeroSlice(fl v.FieldLevel) bool {
	currentValue := fl.Field()
	sliceValue, kind, nullable := fl.ExtractType(currentValue)

	if nullable {
		return true
	}

	switch kind {
	case rf.Array, rf.Slice:
	default:
		return true
	}

	if sliceValue.Len() == 0 {
		return true
	}

	/**
	 * Chooses matching type of element for slice checker
	 */
	var sliceChecker func(rf.Value) bool

	switch sliceValue.Type().Elem().Kind() {
	case rf.Int, rf.Int8, rf.Int16, rf.Int32, rf.Int64:
		sliceChecker = checkSignedInt
	case rf.Uint, rf.Uint8, rf.Uint16, rf.Uint32, rf.Uint64:
		sliceChecker = checkUnsignedInt
	default:
		return true
	}
	// :~)

	return sliceChecker(sliceValue)
}

func checkSignedInt(sliceValue rf.Value) bool {
	for i := 0; i < sliceValue.Len(); i++ {
		if sliceValue.Index(i).Int() != 0 {
			return true
		}
	}

	return false
}
func checkUnsignedInt(sliceValue rf.Value) bool {
	for i := 0; i < sliceValue.Len(); i++ {
		if sliceValue.Index(i).Uint() != 0 {
			return true
		}
	}

	return false
}
