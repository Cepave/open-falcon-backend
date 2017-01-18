package utils

import (
	"fmt"
	"reflect"
)

var (
	TypeOfInt = reflect.TypeOf(int(0))
	TypeOfInt64 = reflect.TypeOf(int64(0))
	TypeOfInt32 = reflect.TypeOf(int32(0))
	TypeOfInt16 = reflect.TypeOf(int16(0))
	TypeOfInt8 = reflect.TypeOf(int8(0))
	TypeOfUint = reflect.TypeOf(uint(0))
	TypeOfUint64 = reflect.TypeOf(uint64(0))
	TypeOfUint32 = reflect.TypeOf(uint32(0))
	TypeOfUint16 = reflect.TypeOf(uint16(0))
	TypeOfUint8 = reflect.TypeOf(uint8(0))

	TypeOfFloat32 = reflect.TypeOf(float32(0))
	TypeOfFloat64 = reflect.TypeOf(float64(0))

	TypeOfComplex64 = reflect.TypeOf(complex64(0))
	TypeOfComplex128 = reflect.TypeOf(complex128(0))

	TypeOfByte = reflect.TypeOf(byte(0))
	TypeOfBool = reflect.TypeOf(true)
	TypeOfString = reflect.TypeOf("")

	ATypeOfInt = reflect.TypeOf([]int{})
	ATypeOfInt64 = reflect.TypeOf([]int64{})
	ATypeOfInt32 = reflect.TypeOf([]int32{})
	ATypeOfInt16 = reflect.TypeOf([]int16{})
	ATypeOfInt8 = reflect.TypeOf([]int8{})
	ATypeOfUint = reflect.TypeOf([]uint{})
	ATypeOfUint64 = reflect.TypeOf([]uint64{})
	ATypeOfUint32 = reflect.TypeOf([]uint32{})
	ATypeOfUint16 = reflect.TypeOf([]uint16{})
	ATypeOfUint8 = reflect.TypeOf([]uint8{})

	ATypeOfFloat32 = reflect.TypeOf([]float32{})
	ATypeOfFloat64 = reflect.TypeOf([]float64{})

	ATypeOfComplex64 = reflect.TypeOf([]complex64{})
	ATypeOfComplex128 = reflect.TypeOf([]complex128{})

	ATypeOfByte = reflect.TypeOf([]byte{})
	ATypeOfBool = reflect.TypeOf([]bool{})
	ATypeOfString = reflect.TypeOf([]string{})

	TrueValue = reflect.ValueOf(true)
	FalseValue = reflect.ValueOf(false)
)

// Super convertion by reflect
//
// 1. Nil pointer would be to nil pointer of target type
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

		if !sourceType.ConvertibleTo(targetType) {
			panic(fmt.Sprintf("Type of source[%s] cannot be converted to type of target[%s]", sourceType, targetType))
		}
	}

	if !sourceType.ConvertibleTo(targetType) {
		panic(fmt.Sprintf("Cannot convert type of [%s] to another type: [%s]", sourceType, targetType))
	}

	return sourceValue.Convert(targetType)
}

type ValueExt reflect.Value
func (v ValueExt) IsArray() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}

	return false
}
func (v ValueExt) IsPointer() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
		return true
	}

	return false
}

// For array, slice, map, chan, the value.Len() must be > 0
//
// For pointer, interface, or function, the value.IsNil() must not be true
//
// Othewise, use value.IsValid()
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
