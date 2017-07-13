package reflect

import (
	"fmt"
	"reflect"
)

var TypeOfReflectType = TypeOfInterface((*reflect.Type)(nil))

func TypeOfValue(v interface{}) reflect.Type {
	if reflectType, ok := v.(reflect.Type); ok {
		return reflectType
	}

	return reflect.TypeOf(v)
}

// Gets the type of interface, which is short function of:
//
// 	reflect.TypeOf(nilPointerToInterface).Elem()
func TypeOfInterface(v interface{}) reflect.Type {
	return reflect.TypeOf(v).Elem()
}

// Checks if the value is type of reflect.Type
func IsReflectType(v interface{}) bool {
	return reflect.TypeOf(v).Implements(TypeOfReflectType)
}

func NewPointerValue(targetValue interface{}) interface{} {
	ptrValue := reflect.New(reflect.TypeOf(targetValue))
	ptrValue.Elem().Set(reflect.ValueOf(targetValue))
	return ptrValue.Interface()
}

// Gets the final pointed type of a pointer type
func FinalPointedType(pointerType reflect.Type) reflect.Type {
	for pointerType.Kind() == reflect.Ptr {
		pointerType = pointerType.Elem()
	}

	return pointerType
}

// Gets the final value of pointer, it may be IsNil()
func FinalPointedValue(pointerValue reflect.Value) reflect.Value {
	for pointerValue.Kind() == reflect.Ptr &&
		!pointerValue.IsNil() {
		pointerValue = pointerValue.Elem()
	}

	return pointerValue
}

// new() a value by pass through multi-layer pointers
func NewFinalValue(startType reflect.Type) reflect.Value {
	newValue := reflect.New(startType).Elem()

	for newValue.Kind() == reflect.Ptr {
		newValue = reflect.New(newValue.Type().Elem()).Elem()
	}

	return newValue
}

func NewFinalValueFrom(fromValue reflect.Value, finalTypeOfPointer reflect.Type) reflect.Value {
	for fromValue.Type() != finalTypeOfPointer {
		oldValue := fromValue

		fromValue = reflect.New(fromValue.Type())
		fromValue.Elem().Set(oldValue)
	}

	return fromValue
}

// Gets the types of a function
//
// This function gives the slice of types for both of input and output
func GetAllTypesForFunction(funcType reflect.Type) (inputTypes []reflect.Type, outputTypes []reflect.Type) {
	if funcType.Kind() != reflect.Func {
		panic(fmt.Sprintf("Need to be function. Got: [%T]", funcType))
	}

	inputTypes = make([]reflect.Type, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		inputTypes[i] = funcType.In(i)
	}

	outputTypes = make([]reflect.Type, funcType.NumOut())
	for i := 0; i < funcType.NumOut(); i++ {
		outputTypes[i] = funcType.Out(i)
	}

	return
}

// Gets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
func GetValueOfField(v interface{}, tree ...string) interface{} {
	return GetValueOfFieldByReflect(
		reflect.ValueOf(v), tree...,
	).Interface()
}

// Sets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
func SetValueOfField(v interface{}, newFieldValue interface{}, tree ...string) {
	SetValueOfFieldByReflect(reflect.ValueOf(v), reflect.ValueOf(newFieldValue), tree...)
}

// Gets value of field, supported tree visiting whether or not the value is
// struct or pointer to struct.
func GetValueOfFieldByReflect(v reflect.Value, tree ...string) reflect.Value {
	currentValue := getValueAsStruct(v)

	for _, fieldName := range tree {
		currentValue = getValueAsStruct(currentValue)
		currentValue = currentValue.FieldByName(fieldName)

		if !currentValue.IsValid() {
			panic(fmt.Sprintf("Field[%s] is INVALID(reflect.Value.IsValid())", fieldName))
		}
	}

	return currentValue
}

// Sets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
func SetValueOfFieldByReflect(v reflect.Value, newFieldValue reflect.Value, tree ...string) {
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("The type of struct must be pointer: [%T]", v))
	}

	GetValueOfFieldByReflect(v, tree...).Set(newFieldValue)
}

func getValueAsStruct(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Needs to be struct. Got: [%v]", v.Type()))
	}

	return v
}
