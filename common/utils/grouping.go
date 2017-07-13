package utils

import (
	"fmt"
	"reflect"
)

type KeyGetter interface {
	GetKey() interface{}
}

type GroupingProcessor struct {
	mapOfKeys     map[interface{}]KeyGetter
	mapOfChildren map[interface{}]reflect.Value

	typeOfElem reflect.Type
}

func NewGroupingProcessor(typeOfElem reflect.Type) *GroupingProcessor {
	return &GroupingProcessor{
		mapOfKeys:     make(map[interface{}]KeyGetter),
		mapOfChildren: make(map[interface{}]reflect.Value),
		typeOfElem:    typeOfElem,
	}
}
func NewGroupingProcessorOfTargetType(target interface{}) *GroupingProcessor {
	return NewGroupingProcessor(reflect.TypeOf(target))
}

func (g *GroupingProcessor) Put(keyObject KeyGetter, child interface{}) {
	keyValue := keyObject.GetKey()

	/**
	 * Puts object of key
	 */
	_, hasKey := g.mapOfKeys[keyValue]
	if !hasKey {
		g.mapOfKeys[keyValue] = keyObject
	}
	// :~)

	/**
	 * Puts object of child
	 */
	children, ok := g.mapOfChildren[keyValue]
	if !ok {
		children = reflect.MakeSlice(
			reflect.SliceOf(g.typeOfElem), 0, 0,
		)
	}

	children = reflect.Append(children, reflect.ValueOf(child))
	g.mapOfChildren[keyValue] = children
	// :~)
}
func (g *GroupingProcessor) Keys() []KeyGetter {
	keys := make([]KeyGetter, len(g.mapOfKeys))

	var i = 0
	for _, keyObject := range g.mapOfKeys {
		keys[i] = keyObject
		i++
	}

	return keys
}
func (g *GroupingProcessor) KeyObject(keyObject KeyGetter) interface{} {
	obj, ok := g.mapOfKeys[keyObject.GetKey()]
	if !ok {
		panic(fmt.Sprintf("Cannot get object for key: [%#v]", keyObject.GetKey()))
	}

	return obj
}
func (g *GroupingProcessor) Children(keyObject KeyGetter) interface{} {
	children, ok := g.mapOfChildren[keyObject.GetKey()]
	if !ok {
		panic(fmt.Sprintf("Cannot get children for key: [%#v]", keyObject.GetKey()))
	}

	return children.Interface()
}
