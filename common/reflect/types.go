package reflect

import (
	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
)

// Hash values of reflect.Type
var (
	/**
	 * Buildin types
	 */
	HashTypeOfInt    = DigestType(TypeOfInt)
	HashTypeOfInt64  = DigestType(TypeOfInt64)
	HashTypeOfInt32  = DigestType(TypeOfInt32)
	HashTypeOfInt16  = DigestType(TypeOfInt16)
	HashTypeOfInt8   = DigestType(TypeOfInt8)
	HashTypeOfUint   = DigestType(TypeOfUint)
	HashTypeOfUint64 = DigestType(TypeOfUint64)
	HashTypeOfUint32 = DigestType(TypeOfUint32)
	HashTypeOfUint16 = DigestType(TypeOfUint16)
	HashTypeOfUint8  = DigestType(TypeOfUint8)

	HashTypeOfFloat32 = DigestType(TypeOfFloat32)
	HashTypeOfFloat64 = DigestType(TypeOfFloat64)

	HashTypeOfComplex64  = DigestType(TypeOfComplex64)
	HashTypeOfComplex128 = DigestType(TypeOfComplex128)

	HashTypeOfByte   = DigestType(TypeOfByte)
	HashTypeOfBool   = DigestType(TypeOfBool)
	HashTypeOfString = DigestType(TypeOfString)
	// :~)

	/**
	 * Slice Types
	 */
	HashPTypeOfInt    = DigestType(PTypeOfInt)
	HashPTypeOfInt64  = DigestType(PTypeOfInt64)
	HashPTypeOfInt32  = DigestType(PTypeOfInt32)
	HashPTypeOfInt16  = DigestType(PTypeOfInt16)
	HashPTypeOfInt8   = DigestType(PTypeOfInt8)
	HashPTypeOfUint   = DigestType(PTypeOfUint)
	HashPTypeOfUint64 = DigestType(PTypeOfUint64)
	HashPTypeOfUint32 = DigestType(PTypeOfUint32)
	HashPTypeOfUint16 = DigestType(PTypeOfUint16)
	HashPTypeOfUint8  = DigestType(PTypeOfUint8)

	HashPTypeOfFloat32 = DigestType(PTypeOfFloat32)
	HashPTypeOfFloat64 = DigestType(PTypeOfFloat64)

	HashPTypeOfComplex64  = DigestType(PTypeOfComplex64)
	HashPTypeOfComplex128 = DigestType(PTypeOfComplex128)

	HashPTypeOfByte   = DigestType(PTypeOfByte)
	HashPTypeOfBool   = DigestType(PTypeOfBool)
	HashPTypeOfString = DigestType(PTypeOfString)
	// :~)

	/**
	 * Slice type
	 */
	HashSTypeOfInt    = DigestType(STypeOfInt)
	HashSTypeOfInt64  = DigestType(STypeOfInt64)
	HashSTypeOfInt32  = DigestType(STypeOfInt32)
	HashSTypeOfInt16  = DigestType(STypeOfInt16)
	HashSTypeOfInt8   = DigestType(STypeOfInt8)
	HashSTypeOfUint   = DigestType(STypeOfUint)
	HashSTypeOfUint64 = DigestType(STypeOfUint64)
	HashSTypeOfUint32 = DigestType(STypeOfUint32)
	HashSTypeOfUint16 = DigestType(STypeOfUint16)
	HashSTypeOfUint8  = DigestType(STypeOfUint8)

	HashSTypeOfFloat32 = DigestType(STypeOfFloat32)
	HashSTypeOfFloat64 = DigestType(STypeOfFloat64)

	HashSTypeOfComplex64  = DigestType(STypeOfComplex64)
	HashSTypeOfComplex128 = DigestType(STypeOfComplex128)

	HashSTypeOfByte   = DigestType(STypeOfByte)
	HashSTypeOfBool   = DigestType(STypeOfBool)
	HashSTypeOfString = DigestType(STypeOfString)
	// :~)
)
