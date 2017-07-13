package types

import (
	"reflect"
)

var (
	/**
	 * Buildin types
	 */
	TypeOfInt    = reflect.TypeOf(int(0))
	TypeOfInt64  = reflect.TypeOf(int64(0))
	TypeOfInt32  = reflect.TypeOf(int32(0))
	TypeOfInt16  = reflect.TypeOf(int16(0))
	TypeOfInt8   = reflect.TypeOf(int8(0))
	TypeOfUint   = reflect.TypeOf(uint(0))
	TypeOfUint64 = reflect.TypeOf(uint64(0))
	TypeOfUint32 = reflect.TypeOf(uint32(0))
	TypeOfUint16 = reflect.TypeOf(uint16(0))
	TypeOfUint8  = reflect.TypeOf(uint8(0))

	TypeOfFloat32 = reflect.TypeOf(float32(0))
	TypeOfFloat64 = reflect.TypeOf(float64(0))

	TypeOfComplex64  = reflect.TypeOf(complex64(0))
	TypeOfComplex128 = reflect.TypeOf(complex128(0))

	TypeOfByte   = reflect.TypeOf(byte(0))
	TypeOfBool   = reflect.TypeOf(true)
	TypeOfString = reflect.TypeOf("")
	// :~)

	/**
	 * Pointer to buildin types
	 */
	PTypeOfInt    = reflect.PtrTo(TypeOfInt)
	PTypeOfInt64  = reflect.PtrTo(TypeOfInt64)
	PTypeOfInt32  = reflect.PtrTo(TypeOfInt32)
	PTypeOfInt16  = reflect.PtrTo(TypeOfInt16)
	PTypeOfInt8   = reflect.PtrTo(TypeOfInt8)
	PTypeOfUint   = reflect.PtrTo(TypeOfUint)
	PTypeOfUint64 = reflect.PtrTo(TypeOfUint64)
	PTypeOfUint32 = reflect.PtrTo(TypeOfUint32)
	PTypeOfUint16 = reflect.PtrTo(TypeOfUint16)
	PTypeOfUint8  = reflect.PtrTo(TypeOfUint8)

	PTypeOfFloat32 = reflect.PtrTo(TypeOfFloat32)
	PTypeOfFloat64 = reflect.PtrTo(TypeOfFloat64)

	PTypeOfComplex64  = reflect.PtrTo(TypeOfComplex64)
	PTypeOfComplex128 = reflect.PtrTo(TypeOfComplex128)

	PTypeOfByte   = reflect.PtrTo(TypeOfByte)
	PTypeOfBool   = reflect.PtrTo(TypeOfBool)
	PTypeOfString = reflect.PtrTo(TypeOfString)
	// :~)

	/**
	 * Slice types of buildin types
	 */
	STypeOfInt    = reflect.TypeOf([]int{})
	STypeOfInt64  = reflect.TypeOf([]int64{})
	STypeOfInt32  = reflect.TypeOf([]int32{})
	STypeOfInt16  = reflect.TypeOf([]int16{})
	STypeOfInt8   = reflect.TypeOf([]int8{})
	STypeOfUint   = reflect.TypeOf([]uint{})
	STypeOfUint64 = reflect.TypeOf([]uint64{})
	STypeOfUint32 = reflect.TypeOf([]uint32{})
	STypeOfUint16 = reflect.TypeOf([]uint16{})
	STypeOfUint8  = reflect.TypeOf([]uint8{})

	STypeOfFloat32 = reflect.TypeOf([]float32{})
	STypeOfFloat64 = reflect.TypeOf([]float64{})

	STypeOfComplex64  = reflect.TypeOf([]complex64{})
	STypeOfComplex128 = reflect.TypeOf([]complex128{})

	STypeOfByte   = reflect.TypeOf([]byte{})
	STypeOfBool   = reflect.TypeOf([]bool{})
	STypeOfString = reflect.TypeOf([]string{})
	// :~)

	TrueValue  = reflect.ValueOf(true)
	FalseValue = reflect.ValueOf(false)
)
