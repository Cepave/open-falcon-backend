package types

import (
	"fmt"
	"reflect"
	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"

	"github.com/Cepave/open-falcon-backend/common/utils"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
)

// Main interface provided by convsersion service
type ConversionService interface {
	CanConvert(sourceType reflect.Type, targetType reflect.Type) bool
	ConvertTo(sourceObject interface{}, targetType reflect.Type) interface{}
}

// Defines the registry for managing converters
type ConverterRegistry interface {
	AddConverter(sourceType reflect.Type, targetType reflect.Type, converter Converter)
}

// Basic conversion function
type Converter func(sourceObject interface{}) interface{}

func NewDefaultConversionService() *DefaultConversionService {
	srv := &DefaultConversionService{
		mapOfConverters: make(map[uint64]map[uint64]Converter),
	}

	return srv
}

// Default implementation for conversion service
//
// First of all, this service would use converters added by AddConverter() function.
//
// If nothing matched, uses buildin converters:
//
//  If target type is pointer, the converted object would be allocated by new().
//
// 	From string source
// 	Convertable to: Int, Uint. Float, Bool
// 		To array: puts the converted element to newArray[0]
// 		To slice: puts the converted element to newSlice[0], len() == 1, cap() == 1
// 		To chan: send the converted element to new(chan <type>, 1)
//
//  From array to slice(convertible element type)
//  From array/slice to channel(convertible element type)
//  From channel to slice(convertible element type)
//
//  From slice to array(panic if: len(slice) > len(array))
//  From channel to array(panic if: len(channel) > len(array))
//
//  From a map to another map of different types:
//  Both of the key and value of maps must be convertible
//
// If the target type is string, this service would use fmt.Sprintf("%v", srcObject) to convert the value by default
//
// Otherwise, see https://golang.org/ref/spec#Conversions
type DefaultConversionService struct {
	mapOfConverters map[uint64]map[uint64]Converter
}
func (s *DefaultConversionService) AddConverter(sourceType reflect.Type, targetType reflect.Type, converter Converter) {
	if sourceType == targetType {
		panic(fmt.Sprintf("Cannot add same type of converter: [%v]", sourceType))
	}

	sourceHash := oreflect.DigestType(sourceType)
	targetHash := oreflect.DigestType(targetType)

	targetFuncs, ok := s.mapOfConverters[sourceHash]
	if !ok {
		targetFuncs = make(map[uint64]Converter)
		s.mapOfConverters[sourceHash] = targetFuncs
	}

	targetFuncs[targetHash] = converter
}
// Only checks added converters(without processing of pointer)
func (s *DefaultConversionService) CanConvert(sourceType reflect.Type, targetType reflect.Type) bool {
	sourceHash := oreflect.DigestType(sourceType)
	targetHash := oreflect.DigestType(targetType)

	_, ok := s.mapOfConverters[sourceHash][targetHash]
	return ok
}
// Converts data by converter
//
// If the target type is pointer to type, which could be converted, this function would use it automatically.
func (s *DefaultConversionService) ConvertTo(sourceObject interface{}, targetType reflect.Type) interface{} {
	sourceType := reflect.TypeOf(sourceObject)
	if sourceType == targetType {
		return sourceObject
	}

	sourceHash := oreflect.DigestType(sourceType)
	targetHash := oreflect.DigestType(targetType)

	targetFunc, ok := s.mapOfConverters[sourceHash][targetHash]
	if ok {
		return targetFunc(sourceObject)
	}

	return s.defaultConvert(
		sourceObject, targetType,
		sourceHash, targetHash,
	)
}

func (s *DefaultConversionService) convertToReflect(sourceObject interface{}, targetType reflect.Type) reflect.Value {
	return reflect.ValueOf(s.ConvertTo(sourceObject, targetType))
}
func (s *DefaultConversionService) convertFromReflectToReflect(sourceValue reflect.Value, targetType reflect.Type) reflect.Value {
	return reflect.ValueOf(s.ConvertTo(sourceValue.Interface(), targetType))
}

func (s *DefaultConversionService) defaultConvert(
	sourceObject interface{}, targetType reflect.Type,
	sourceHash uint64, targetHash uint64,
) interface{} {
	/**
	 * // Converts any type to string
	 */
	if targetType.Kind() == reflect.String {
		return fmt.Sprintf("%v", sourceObject)
	}
	// :~)

	if reflect.TypeOf(sourceObject).ConvertibleTo(targetType) {
		return reflect.ValueOf(sourceObject).Convert(targetType).Interface()
	}

	/**
	 * Converts the source(pointer to something) to final value
	 */
	if reflect.TypeOf(sourceObject).Kind() == reflect.Ptr {
		sourceObjectV := oreflect.FinalPointedValue(reflect.ValueOf(sourceObject))
		/**
		 * The pointed value is nil pointer
		 */
		if sourceObjectV.Kind() == reflect.Ptr {
			panic(fmt.Sprintf("Convert nil value to type[%v]. Nil: [%v]", targetType, sourceObjectV.Type()))
		}
		// :~)

		sourceObject = sourceObjectV.Interface()
		sourceHash = oreflect.DigestType(sourceObjectV.Type())
	}
	// :~)

	switch targetType.Kind() {
	case reflect.Ptr:
		return s.convertToTypeOfPtr(sourceObject, targetType)
	case reflect.Chan:
		return s.convertToTypeOfChan(sourceObject, targetType)
	case reflect.Slice:
		return s.convertToTypeOfSlice(sourceObject, targetType)
	case reflect.Array:
		return s.convertToTypeOfArray(sourceObject, targetType)
	case reflect.Map:
		if reflect.TypeOf(sourceObject).Kind() == reflect.Map {
			return s.convertToAnotherMap(sourceObject, targetType)
		}
	}

	converter, ok := buildinConverters[sourceHash][targetHash]
	if ok {
		return converter(sourceObject)
	}

	// Uses default convertion of types
	return utils.ConvertTo(sourceObject, targetType)
}

func (s *DefaultConversionService) convertToTypeOfSlice(srcObject interface{}, targetType reflect.Type) interface{} {
	targetTypeOfElem := targetType.Elem()
	valueOfSrc := reflect.ValueOf(srcObject)
	srcKind := valueOfSrc.Kind()
	switch srcKind {
	case reflect.Slice, reflect.Array, reflect.Chan:
		container := reflect.MakeSlice(targetType, 0, valueOfSrc.Len())
		containerWalkers[srcKind](
			valueOfSrc,
			func(index int, obj reflect.Value) {
				container = reflect.Append(
					container,
					s.convertFromReflectToReflect(obj, targetTypeOfElem),
				)
			},
		)
		return container.Interface()
	}

	sliceValue := reflect.MakeSlice(targetType, 1, 1)
	sliceValue.Index(0).Set(
		s.convertToReflect(srcObject, targetTypeOfElem),
	)

	return sliceValue.Interface()
}
func (s *DefaultConversionService) convertToTypeOfArray(srcObject interface{}, targetType reflect.Type) interface{} {
	targetTypeOfElem := targetType.Elem()
	valueOfSrc := reflect.ValueOf(srcObject)
	srcKind := valueOfSrc.Kind()
	switch srcKind {
	case reflect.Slice, reflect.Array, reflect.Chan:
		if valueOfSrc.Len() > targetType.Len() {
			panic(fmt.Sprintf(
				"Cannot convert type[%T] to type[%s]. Len of source: [%d]",
				srcObject, targetType, valueOfSrc.Len(),
			))
		}
		container := reflect.New(targetType).Elem()
		containerWalkers[srcKind](
			valueOfSrc,
			func(index int, obj reflect.Value) {
				container.Index(index).Set(
					s.convertFromReflectToReflect(obj, targetTypeOfElem),
				)
			},
		)

		return container.Interface()
	}

	arrayValue := reflect.New(targetType).Elem()
	arrayValue.Index(0).Set(
		s.convertToReflect(srcObject, targetTypeOfElem),
	)
	return arrayValue.Interface()
}
func (s *DefaultConversionService) convertToTypeOfPtr(srcObject interface{}, targetType reflect.Type) interface{} {
	targetType = targetType.Elem()

	pointerResult := reflect.New(targetType)
	pointerResult.Elem().Set(
		s.convertToReflect(srcObject, targetType),
	)

	return pointerResult.Interface()
}
func (s *DefaultConversionService) convertToTypeOfChan(srcObject interface{}, targetType reflect.Type) interface{} {
	targetTypeOfElem := targetType.Elem()
	valueOfSrc := reflect.ValueOf(srcObject)
	srcKind := valueOfSrc.Kind()
	switch srcKind {
	case reflect.Slice, reflect.Array, reflect.Chan:
		channelValue := reflect.MakeChan(targetType, valueOfSrc.Len())

		containerWalkers[srcKind](
			valueOfSrc,
			func(index int, obj reflect.Value) {
				obj = s.convertFromReflectToReflect(obj, targetTypeOfElem)
				if !channelValue.TrySend(obj) {
					panic(fmt.Sprintf(
						"Cannot send value to channel. Type:[%s]. Object: [%#v]",
						targetType, obj.Interface(),
					))
				}
			},
		)
		return channelValue.Interface()
	}

	channelValue := reflect.MakeChan(targetType, 1)
	channelValue.Send(
		s.convertToReflect(srcObject, targetTypeOfElem),
	)

	return channelValue.Interface()
}
func (s *DefaultConversionService) convertToAnotherMap(srcObject interface{}, targetType reflect.Type) interface{} {
	srcValue := reflect.ValueOf(srcObject)
	newMap := reflect.MakeMap(targetType)

	for _, key := range srcValue.MapKeys() {
		convertedKey := s.convertFromReflectToReflect(key, targetType.Key())
		convertedValue := s.convertFromReflectToReflect(srcValue.MapIndex(key), targetType.Elem())

		newMap.SetMapIndex(convertedKey, convertedValue)
	}

	return newMap.Interface()
}

var buildinConverters = map[uint64]map[uint64]Converter {
	oreflect.HashTypeOfString: map[uint64]Converter {
		oreflect.HashTypeOfInt: func(v interface{}) interface{} { return convertStringToInt(v, TypeOfInt) },
		oreflect.HashTypeOfInt64: func(v interface{}) interface{} { return convertStringToInt(v, TypeOfInt64) },
		oreflect.HashTypeOfInt32: func(v interface{}) interface{} { return convertStringToInt(v, TypeOfInt32) },
		oreflect.HashTypeOfInt16: func(v interface{}) interface{} { return convertStringToInt(v, TypeOfInt16) },
		oreflect.HashTypeOfInt8: func(v interface{}) interface{} { return convertStringToInt(v, TypeOfInt8) },
		oreflect.HashTypeOfUint: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfUint) },
		oreflect.HashTypeOfUint64: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfUint64) },
		oreflect.HashTypeOfUint32: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfUint32) },
		oreflect.HashTypeOfUint16: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfUint16) },
		oreflect.HashTypeOfUint8: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfUint8) },

		oreflect.HashTypeOfFloat32: func(v interface{}) interface{} { return convertStringToFloat(v, TypeOfFloat32) },
		oreflect.HashTypeOfFloat64: func(v interface{}) interface{} { return convertStringToFloat(v, TypeOfFloat64) },

		oreflect.HashTypeOfByte: func(v interface{}) interface{} { return convertStringToUint(v, TypeOfByte) },
		oreflect.HashTypeOfBool: func(v interface{}) interface{} { return convertStringToBool(v, TypeOfBool) },
	},
}

type containerWalker func(reflect.Value, func(int, reflect.Value))
var containerWalkers = map[reflect.Kind]containerWalker {
	reflect.Slice: func(src reflect.Value, callback func(int, reflect.Value)) {
		for i := 0; i < src.Len(); i++ {
			callback(i, src.Index(i))
		}
	},
	reflect.Array: func(src reflect.Value, callback func(int, reflect.Value)) {
		for i := 0; i < src.Len(); i++ {
			callback(i, src.Index(i))
		}
	},
	reflect.Chan: func(src reflect.Value, callback func(int, reflect.Value)) {
		channelLen := src.Len()
		for i := 0; i < channelLen; i++ {
			obj, ok := src.TryRecv()
			if !ok {
				panic(fmt.Sprintf("Cannot receive from channel. Channel type:[%v]", src.Type()))
			}

			callback(i, obj)
			src.Send(obj) // Put the object back to channel
		}
	},
}
