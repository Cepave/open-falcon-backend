package types

import (
	"fmt"
	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
	"reflect"
)

// Defines the binding interface to convert any object to implementing type
type Binding interface {
	// Binds the content of source object into this object
	Bind(sourceObject interface{})
}

// Convenient function to perform the binding
//
// 	sourceObject - The source to be converted
// 	holder - The object implements "Binding" interface
//
// This function would panic if the holder doesn't implement "Binding" interface.
func DoBinding(sourceObject interface{}, holder interface{}) {
	b, ok := holder.(Binding)
	if !ok {
		panic(fmt.Sprintf("Type: [%t] doesn't implmenet \"Binding\" interface", holder))
	}

	b.Bind(sourceObject)
}

// Checks if a object has implemented "Binding" interface
func HasBinding(holder interface{}) bool {
	_, ok := holder.(Binding)
	return ok
}

// Converts the Binding(with builder) to "Converter"
func BindingToConverter(holderBuilder func() interface{}) (Converter, reflect.Type) {
	converter := func(object interface{}) interface{} {
		o := holderBuilder()
		DoBinding(object, o)
		return o
	}

	return converter, reflect.TypeOf(holderBuilder())
}

var _t_Binding = oreflect.TypeOfInterface((*Binding)(nil))
