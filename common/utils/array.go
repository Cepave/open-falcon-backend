package utils

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	TrimStringMapper = TypedFuncToMapper(func(v string) string {
		return strings.TrimSpace(v)
	})

	EmptyStringFilter = TypedFuncToFilter(func(v string) bool {
		return v != ""
	})
)

// This function is used to get bool value to decide
type FilterFunc func(interface{}) bool

// This function is used to transfer element
type MapperFunc func(interface{}) interface{}

func IdentityMapper(v interface{}) interface{} {
	return v
}

// Converts typed function(for filter) to FilterFunc
func TypedFuncToFilter(anyFunc interface{}) FilterFunc {
	valueOfFunc := reflect.ValueOf(anyFunc)
	typeOfFunc := valueOfFunc.Type()

	if typeOfFunc.NumIn() != 1 || typeOfFunc.NumIn() != 1 ||
		typeOfFunc.Out(0) != TypeOfBool {
		panic(fmt.Errorf("Filter need to be type of \"func(interface{}) bool\""))
	}

	return func(v interface{}) bool {
		funcInputType := typeOfFunc.In(0)

		inputValue := convertGenericValue(v, funcInputType)
		outputValue := valueOfFunc.Call([]reflect.Value{ inputValue })[0]

		return outputValue.Bool()
	}
}

// Constructs a filter, which uses map of golang.
//
// **WARNING** The generated filter cannot be reused
//
// 	targetType - Could be array, slice, or element type
func NewUniqueFilter(targetType reflect.Type) FilterFunc {
	var typeOfMap reflect.Type

	switch targetType.Kind() {
	case reflect.Array, reflect.Slice:
		typeOfEle := targetType.Elem()
		typeOfMap = reflect.MapOf(typeOfEle, TypeOfBool)
	default:
		typeOfMap = reflect.MapOf(targetType, TypeOfBool)
	}

	valueOfUniqueMap := reflect.MakeMap(typeOfMap)

	return func(v interface{}) bool {
		key := reflect.ValueOf(v)

		uniqueValue := valueOfUniqueMap.MapIndex(key)
		if uniqueValue.IsValid() {
			return false
		}

		valueOfUniqueMap.SetMapIndex(key, TrueValue)
		return true
	}
}

// Constructs a filter, which uses domain(map[<type]bool) of golang as filtering
//
// The element must be shown in domain.
func NewDomainFilter(mapOfDomain interface{}) FilterFunc {
	valueOfDomain := reflect.ValueOf(mapOfDomain)

	return func(v interface{}) bool {
		value := valueOfDomain.MapIndex(reflect.ValueOf(v))
		return value.IsValid()
	}
}

// Converts typed function(for filter) to MapperFunc
func TypedFuncToMapper(anyFunc interface{}) MapperFunc {
	mapperTypes := getMapperTypes(anyFunc)

	funcValue := reflect.ValueOf(anyFunc)
	return func(v interface{}) interface{} {
		inputValue := convertGenericValue(v, mapperTypes[0])
		outputValue := funcValue.Call([]reflect.Value{ inputValue })[0]

		return outputValue.Interface()
	}
}

// Abstract array to provide various function for processing array
type AbstractArray struct {
	arrayElementType reflect.Type
	anyArrayValue reflect.Value
}

// Constructs an array of abstract
func MakeAbstractArray(sourceSlice interface{}) *AbstractArray {
	valueOfArray := reflect.ValueOf(sourceSlice)

	switch valueOfArray.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		panic(fmt.Errorf("Cannot support of type for abstract array: %v", valueOfArray.Kind()))
	}

	return &AbstractArray { valueOfArray.Type().Elem(), valueOfArray }
}

// Gets the result array
func (a *AbstractArray) GetArray() interface{} {
	return a.GetArrayAsType(a.arrayElementType)
}

func (a *AbstractArray) GetArrayAsType(targetType reflect.Type) interface{} {
	valueOfAnyArray := a.anyArrayValue
	size := valueOfAnyArray.Len()

	var newArrayType reflect.Type
	var convertValue func(srcValue reflect.Value) reflect.Value
	if a.arrayElementType == targetType {
		newArrayType = valueOfAnyArray.Type()
		convertValue = func(srcValue reflect.Value) reflect.Value { return srcValue }
	} else {
		newArrayType = reflect.SliceOf(targetType)
		convertValue = func(srcValue reflect.Value) reflect.Value {
			return convertAnyValue(srcValue, targetType)
		}
	}

	newArrayValue := reflect.MakeSlice(newArrayType, size, size)

	for i := 0; i < size; i++ {
		newArrayValue.Index(i).Set(
			convertValue(valueOfAnyArray.Index(i)),
		)
	}

	return newArrayValue.Interface()
}

// Filters elements in the array
func (a *AbstractArray) FilterWith(filter FilterFunc) *AbstractArray {
	valueOfAnyArray := a.anyArrayValue

	newArray := reflect.MakeSlice(valueOfAnyArray.Type(), 0, 0)

	for i := 0; i < valueOfAnyArray.Len(); i++ {
		currentValue := valueOfAnyArray.Index(i)

		if !filter(currentValue.Interface()) {
			continue
		}

		newArray = reflect.Append(newArray, currentValue)
	}

	return MakeAbstractArray(newArray.Interface())
}
// Maps the elements in array(with target type of result array)
func (a *AbstractArray) MapTo(mapper MapperFunc, eleType reflect.Type) *AbstractArray {
	valueOfAnyArray := a.anyArrayValue

	newArray := reflect.MakeSlice(
		reflect.SliceOf(eleType),
		valueOfAnyArray.Len(), valueOfAnyArray.Len(),
	)
	// :~)

	for i := 0; i < valueOfAnyArray.Len(); i++ {
		currentValue := valueOfAnyArray.Index(i)

		transferedValue := reflect.ValueOf(mapper(currentValue.Interface()))
		transferedValue = convertAnyValue(transferedValue, eleType)

		newArray.Index(i).Set(transferedValue)
	}

	return MakeAbstractArray(newArray.Interface())
}

func getMapperTypes(mapperFunc interface{}) []reflect.Type {
	funcType := reflect.TypeOf(mapperFunc)

	if funcType.NumIn() != 1 || funcType.NumIn() != 1 {
		panic(fmt.Errorf("Need in(1) and out(1) for mapper func"))
	}

	return []reflect.Type{ funcType.In(0), funcType.Out(0) }
}
func convertAnyValue(value reflect.Value, targetType reflect.Type) reflect.Value {
	if value.Type() == targetType {
		return value
	}

	if !value.Type().ConvertibleTo(targetType) {
		panic(fmt.Errorf("Cannot convert type of [%v] to type of [%v]", value.Type(), targetType))
	}

	return value.Convert(targetType)
}

func convertGenericValue(v interface{}, targetType reflect.Type) reflect.Value {
	return convertAnyValue(reflect.ValueOf(v), targetType)
}
