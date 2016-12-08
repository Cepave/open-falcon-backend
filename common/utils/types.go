package utils

import (
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

	TypeOfByte = reflect.TypeOf(byte(0))
	TypeOfBool = reflect.TypeOf(true)
	TypeOfString = reflect.TypeOf("")

	TrueValue = reflect.ValueOf(true)
	FalseValue = reflect.ValueOf(false)
)
