package utils

import (
	"fmt"
	"reflect"
)

// TODO 以下的部分, 考虑放到公共组件库
func KeysOfMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}

type AbstractMap struct {
	currentMap interface{}
}

func MakeAbstractMap(sourceMap interface{}) *AbstractMap {
	valueOfMap := reflect.ValueOf(sourceMap)

	if valueOfMap.Kind() != reflect.Map {
		panic(fmt.Sprintf("The type of object is not map: [%T]", sourceMap))
	}

	return &AbstractMap{sourceMap}
}

// Gets map of desired type
func (m *AbstractMap) ToType(keyType reflect.Type, elemType reflect.Type) interface{} {
	valueOfSourceMap := reflect.ValueOf(m.currentMap)

	resultMap := reflect.MakeMap(reflect.MapOf(keyType, elemType))

	for _, key := range valueOfSourceMap.MapKeys() {
		sourceElem := valueOfSourceMap.MapIndex(key)

		targetKey := ConvertToByReflect(key, keyType)
		targetElem := ConvertToByReflect(sourceElem, elemType)

		resultMap.SetMapIndex(targetKey, targetElem)
	}

	return resultMap.Interface()
}

func (m *AbstractMap) ToTypeOfTarget(keyOfTarget interface{}, elemOfTarget interface{}) interface{} {
	return m.ToType(
		reflect.TypeOf(keyOfTarget),
		reflect.TypeOf(elemOfTarget),
	)
}
