package reflect

import (
	"fmt"
	"reflect"
)

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
			panic(fmt.Sprintf("Field[%s] is zero(reflect.Value.IsZero())", fieldName))
		}
	}

	return currentValue
}
// Sets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
func SetValueOfFieldByReflect(v reflect.Value, newFieldValue reflect.Value, tree ...string) {
	GetValueOfFieldByReflect(v, tree...).Set(newFieldValue)
}

func getValueAsStruct(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Needs to be struct. Got: [%v]", v.Type()))
	}

	return v
}
