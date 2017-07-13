package digest

import (
	"bytes"
	"encoding/binary"
	"fmt"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var Logger = log.NewDefaultLogger("WARN")

type BytesGetter func() []byte

const DigestTagName = "digest"

var zeroBytes = []byte{0}

// Gets the function provides constant zero bytes(not empty array)
func ZeroBytesFunc() []byte {
	return zeroBytes
}

// Gets the function provides constant non-zero bytes
var nonZeroBytes = []byte{1}

func NonZeroBytesFunc() []byte {
	return nonZeroBytes
}

var noByte = []byte{}

// Gets the function provides constant no byte(empty array)
func NoByteFunc() []byte {
	return noByte
}

// This package supports the digesting on struct:
//
// 		type MyStruct type {
// 			Name string `digest:1`
// 			Age int `digest:2`
// 		}
//
// Constructs the struct tag and performs digesting
//
// Digesting rule on field: []byte(<field name> +'|') bytes on fields
//
// If field is the type of "interface{}",
// this type would be purified into more accurate type, which makes the digesting more easy.
//
// Supporting type on field:
// 	buildin types of golang
// 	pointers of buildin types of golang
//
// 	array types(element could be any type or pointer types)
//
// 	the type implementing Digestor
//
// 	"interface{}"
//
// Fields could implement Digestor for customized bytes of your object:
//
//		type sha1String string
//		func (s sha1String) GetDigest() []byte {
//			digest := sha1.Sum([]byte(s))
//			return digest[:]
//		}
//
// According to the "sequence" part of tag value,
// the digest value of every field would be concated to a stream of bytes and
// use sumFunc to digest the final result.
//
// For nested structs, the inner one would be digested(by sumFunc) as stream of bytes.
//
// For bool values:
//
// 	true - []byte{ 1 }
// 	false - []byte{ 0 }
func DigestStruct(input interface{}, sumFunc SumFunc) []byte {
	if input == nil {
		panic("The instance of struct is nil")
	}

	valueOfInput := reflect.ValueOf(input)
	typeOfInput := valueOfInput.Type()

	/**
	 * Checks the type of struct
	 */
	switch typeOfInput.Kind() {
	case reflect.Ptr:
		extractedValue := extractUntilNilOrNotPointer(valueOfInput)
		return DigestStruct(extractedValue.Interface(), sumFunc)
	case reflect.Struct:
		return digestStructForReal(valueOfInput, sumFunc)
	}
	// :~)

	panic(fmt.Errorf("Need to be struct. Got: %s", typeOfInput.Name()))
}

// Gets the bytes getter by any value,
// this function could be used to implement your own DigestStruct
func GetBytesGetter(v interface{}, sumFunc SumFunc) BytesGetter {
	for {
		switch v.(type) {
		case reflect.Value:
			v = v.(reflect.Value).Interface()
		default:
			return buildBytesGetter(reflect.ValueOf(v), sumFunc)
		}
	}
}
func buildBytesGetter(value reflect.Value, sumFunc SumFunc) BytesGetter {
	checkedType := value.Type()

	valueAsGeneric := value.Interface()
	digestor, isDigestor := valueAsGeneric.(Digestor)
	if isDigestor {
		Logger.Debugf("The object of type [%s] is Digestor.", checkedType)
		return digestor.GetDigest
	}

	Logger.Debugf("%s", &showType{checkedType})

	switch checkedType.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr:
		return buildPointerBytesGetter(value, sumFunc)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return buildBytesGetterByBigEndianBinaryForValue(value)
	case reflect.Int:
		return buildBytesGetterByBigEndianBinary(value.Int())
	case reflect.Uint:
		return buildBytesGetterByBigEndianBinary(value.Uint())
	case reflect.Bool:
		return buildBoolBytesGetter(value.Bool())
	case reflect.Array, reflect.Slice:
		return buildArrayBytesGetter(value, sumFunc)
	case reflect.Struct:
		return func() []byte {
			return digestStructForReal(value, sumFunc)
		}
	case reflect.String:
		return buildStringBytesGetter(value)
	case reflect.Interface:
		return GetBytesGetter(reflect.ValueOf(valueAsGeneric), sumFunc)
	}

	panic(fmt.Errorf("Unsupported type for digesting: %v", checkedType))
}

func digestStructForReal(valueOfInput reflect.Value, sumFunc SumFunc) []byte {
	typeOfInput := valueOfInput.Type()

	/**
	 * Loads the ordered fields to be digested
	 */
	orderedFields := loadDigestingSequence(typeOfInput)
	if len(orderedFields) == 0 {
		panic(fmt.Errorf("Nothing to digest for struct: %s", typeOfInput.Name()))
	}
	// :~)

	Logger.Debugf("Found [%d] fields to be digested. For type: [%s]", len(orderedFields), typeOfInput)

	/**
	 * Digests the value of fields in sequence
	 */
	bytesOfStruct := make([]byte, 0)
	for _, fieldToBeDigested := range orderedFields {
		fieldValue := valueOfInput.FieldByName(fieldToBeDigested.name)

		Logger.Debugf("Field Name: [%s].", fieldToBeDigested.name)

		funcBytesGetter := GetBytesGetter(fieldValue, sumFunc)
		bytesOfField := funcBytesGetter()

		if len(bytesOfField) == 0 {
			continue
		}

		bytesOfStruct = append(
			bytesOfStruct, []byte(fieldToBeDigested.name+"|")...,
		)
		bytesOfStruct = append(
			bytesOfStruct, bytesOfField...,
		)
	}
	// :~)

	return sumFunc(bytesOfStruct)
}

func buildPointerBytesGetter(v reflect.Value, sumFunc SumFunc) BytesGetter {
	extractedValue := extractUntilNilOrNotPointer(v)

	if extractedValue.Type().Kind() == reflect.Ptr {
		if extractedValue.IsNil() {
			Logger.Debugf("Nil pointer. Use []byte{} as data for digesting")
			return NoByteFunc
		}
	}

	return GetBytesGetter(extractedValue, sumFunc)
}

var typeOfString = reflect.TypeOf("")

func buildStringBytesGetter(v reflect.Value) BytesGetter {
	stringValue := v.Convert(typeOfString)

	return func() []byte {
		return []byte(stringValue.Interface().(string))
	}
}
func buildBytesGetterByBigEndianBinaryForValue(value reflect.Value) BytesGetter {
	return buildBytesGetterByBigEndianBinary(value.Interface())
}
func buildBytesGetterByBigEndianBinary(value interface{}) BytesGetter {
	return func() []byte {
		return toBigEndianBytes(value)
	}
}
func toBigEndianBytes(v interface{}) []byte {
	var byteBuffer bytes.Buffer

	err := binary.Write(&byteBuffer, binary.BigEndian, v)
	if err != nil {
		panic(fmt.Errorf(
			"Cannot write data as big endian(encoding/binary): \"%v\". Error: %v",
			v, err,
		))
	}

	return byteBuffer.Bytes()
}
func buildBoolBytesGetter(boolValue bool) BytesGetter {
	if boolValue {
		return NonZeroBytesFunc
	}

	return ZeroBytesFunc
}
func digestorToBytes(v interface{}) []byte {
	return v.(Digestor).GetDigest()
}

func buildArrayBytesGetter(
	arrayValue reflect.Value, sumFunc SumFunc,
) BytesGetter {
	bytesValue, isBytesValue := arrayValue.Interface().([]byte)
	if isBytesValue {
		Logger.Debugf("Output array of bytes. Len: %d", len(bytesValue))
		return func() []byte {
			return bytesValue
		}
	}

	if arrayValue.Len() == 0 {
		Logger.Debugf("Array size is zero")
		return NoByteFunc
	}

	return func() []byte {
		allBytes := make([]byte, 0)

		lenOfArray := arrayValue.Len()
		for i := 0; i < lenOfArray; i++ {
			Logger.Debugf("Process element of array [%d]", i)
			bytesFunc := GetBytesGetter(arrayValue.Index(i), sumFunc)
			bytesOfElement := bytesFunc()

			allBytes = append(allBytes, bytesOfElement...)
		}

		return allBytes
	}
}

type fieldWithSequence struct {
	name     string
	sequence int
}

func (f *fieldWithSequence) String() string {
	return fmt.Sprintf("|%d| Name: [%s].", f.sequence, f.name)
}

type sortingFields []*fieldWithSequence

func (f sortingFields) Len() int           { return len(f) }
func (f sortingFields) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f sortingFields) Less(i, j int) bool { return f[i].sequence < f[j].sequence }

func loadDigestingSequence(typeOfStruct reflect.Type) sortingFields {
	fieldsToBeDigested := make(sortingFields, 0)

	for i := 0; i < typeOfStruct.NumField(); i++ {
		field := typeOfStruct.Field(i)

		if !unicode.IsUpper([]rune(field.Name)[0]) {
			continue
		}

		tagValue := field.Tag.Get(DigestTagName)
		if tagValue == "" {
			continue
		}

		tagData := strings.Split(tagValue, ",")

		if len(tagData) != 1 {
			panic(fmt.Errorf("Cannot parse tag value for digest: %s", tagValue))
		}

		/**
		 * Process the sequence of digesting
		 */
		sequence, err := strconv.ParseInt(tagData[0], 10, 32)
		if err != nil {
			panic(fmt.Errorf("Cannot parse tag value for sequence: %s", tagData[0]))
		}
		// :~)

		fieldsToBeDigested = append(
			fieldsToBeDigested,
			&fieldWithSequence{field.Name, int(sequence)},
		)
	}

	sort.Sort(fieldsToBeDigested)
	return fieldsToBeDigested
}

type showType struct {
	typeObject reflect.Type
}

func (t *showType) String() string {
	typeObject := t.typeObject

	return fmt.Sprintf(
		"Type name: [%s]. Kind: [%s]",
		typeObject.Name(), typeObject.Kind(),
	)
}

func extractUntilNilOrNotPointer(value reflect.Value) reflect.Value {
	for {
		typeOfPointer := value.Type()

		if typeOfPointer.Kind() != reflect.Ptr {
			return value
		}

		if value.IsNil() {
			return value
		}

		value = value.Elem()
	}
}
