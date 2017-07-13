package utils

import (
	"fmt"
	"reflect"

	ot "github.com/Cepave/open-falcon-backend/common/reflect/types"
)

var (
	TypeOfInt    = ot.TypeOfInt
	TypeOfInt64  = ot.TypeOfInt64
	TypeOfInt32  = ot.TypeOfInt32
	TypeOfInt16  = ot.TypeOfInt16
	TypeOfInt8   = ot.TypeOfInt8
	TypeOfUint   = ot.TypeOfUint
	TypeOfUint64 = ot.TypeOfUint64
	TypeOfUint32 = ot.TypeOfUint32
	TypeOfUint16 = ot.TypeOfUint16
	TypeOfUint8  = ot.TypeOfUint8

	TypeOfFloat32 = ot.TypeOfFloat32
	TypeOfFloat64 = ot.TypeOfFloat64

	TypeOfComplex64  = ot.TypeOfComplex64
	TypeOfComplex128 = ot.TypeOfComplex128

	TypeOfByte   = ot.TypeOfByte
	TypeOfBool   = ot.TypeOfBool
	TypeOfString = ot.TypeOfString

	ATypeOfInt    = ot.STypeOfInt
	ATypeOfInt64  = ot.STypeOfInt64
	ATypeOfInt32  = ot.STypeOfInt32
	ATypeOfInt16  = ot.STypeOfInt16
	ATypeOfInt8   = ot.STypeOfInt8
	ATypeOfUint   = ot.STypeOfUint
	ATypeOfUint64 = ot.STypeOfUint64
	ATypeOfUint32 = ot.STypeOfUint32
	ATypeOfUint16 = ot.STypeOfUint16
	ATypeOfUint8  = ot.STypeOfUint8

	ATypeOfFloat32 = ot.STypeOfFloat32
	ATypeOfFloat64 = ot.STypeOfFloat64

	ATypeOfComplex64  = ot.STypeOfComplex64
	ATypeOfComplex128 = ot.STypeOfComplex128

	ATypeOfByte   = ot.STypeOfByte
	ATypeOfBool   = ot.STypeOfBool
	ATypeOfString = ot.STypeOfString

	TrueValue  = ot.TrueValue
	FalseValue = ot.FalseValue
)

// If the value of source is empty(""), gets nil pointer
func PointerOfCloneString(source string) *string {
	if source == "" {
		return nil
	}

	return &source
}

// Super conversion by reflect
//
// 1. Nil pointer would be to nil pointer of target type
// 2. Othewise, uses the reflect.Value.Convert() function to perform conversion
//
// See https://golang.org/ref/spec#Conversions
func ConvertTo(value interface{}, targetType reflect.Type) interface{} {
	return ConvertToByReflect(
		reflect.ValueOf(value), targetType,
	).Interface()
}
func ConvertToTargetType(value interface{}, targetValue interface{}) interface{} {
	return ConvertTo(value, reflect.TypeOf(targetValue))
}
func ConvertToByReflect(sourceValue reflect.Value, targetType reflect.Type) reflect.Value {
	sourceType := sourceValue.Type()

	if sourceType == targetType {
		return sourceValue
	}

	switch sourceType.Kind() {
	case reflect.Ptr:
		if targetType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("Target type is not pointer:[%s]. Source type:[%s]", targetType, sourceType))
		}

		if sourceValue.IsNil() {
			return reflect.Zero(targetType)
		}
	}

	if !sourceType.ConvertibleTo(targetType) {
		panic(fmt.Sprintf("Cannot convert type of [%s] to another type: [%s]", sourceType, targetType))
	}

	return sourceValue.Convert(targetType)
}

func IsViable(v interface{}) bool {
	return ValueExt(reflect.ValueOf(v)).IsViable()
}

// Alias of reflect.Value, provides some convenient functions for programming on reflection.
type ValueExt reflect.Value

// Returns true value if the value is array or slice
func (v ValueExt) IsArray() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}

	return false
}

// Returns true value if the value is reflect.Ptr, reflect.Uintptr, or reflect.UnsafePointer
func (v ValueExt) IsPointer() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
		return true
	}

	return false
}

// Checks if a value is viable
//
// 	For array, slice, map, chan: the value.Len() must be > 0
//
//	For pointer, interface, or function: the value.IsNil() must not be true
//
//	Othewise: use reflect.Value.IsValid()
func (v ValueExt) IsViable() bool {
	reflectValue := reflect.Value(v)

	switch reflectValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return reflectValue.Len() > 0
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Func:
		return !reflectValue.IsNil()
	default:
		return reflectValue.IsValid()
	}
}
