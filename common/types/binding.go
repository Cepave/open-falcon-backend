package types

import (
	"fmt"
	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
	"reflect"
)

// Performs the binding
func DoBinding(sourceObject interface{}, holder interface{}) {
	b, ok := holder.(Binding)
	if !ok {
		panic(fmt.Sprintf("Type: [%t] doesn't implmenet \"Binding\" interface", holder))
	}

	b.Bind(sourceObject)
}
// Checks if the interface has implmented "Binding" interface
func HasBinding(holder interface{}) bool {
	_, ok := holder.(Binding)
	return ok
}

// Converts the Binding(with builder) to Converter
func BindingToConverter(holderBuilder func() interface{}) (Converter, reflect.Type) {
	converter := func(object interface{}) interface{} {
		o := holderBuilder()
		DoBinding(object, o)
		return o
	}

	return converter, reflect.TypeOf(holderBuilder).Out(0)
}

// Defines the binding interface to convert any object to implementing type
type Binding interface {
	Bind(sourceObject interface{})
}

var _t_Binding = oreflect.TypeOfInterface((*Binding)(nil))
