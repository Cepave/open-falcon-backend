package textbuilder

import (
	"fmt"
)

// Build-in double-quoted transformers
var QQ = map[string]Transformer {
	"()" : BuildSurrounding(Dsl.S("("), Dsl.S(")")),
	"[]" : BuildSurrounding(Dsl.S("["), Dsl.S("]")),
	"{}" : BuildSurrounding(Dsl.S("{"), Dsl.S("}")),
	"<>" : BuildSurrounding(Dsl.S("<"), Dsl.S(">")),
	"\"" : BuildSameSurrounding(Dsl.S("\"")),
	"'" : BuildSameSurrounding(Dsl.S("'")),
}

var J = map[string]Distiller {
	"," : BuildJoin(Dsl.S(",")),
	", " : BuildJoin(Dsl.S(", ")),
}

// Short name of building blocks for text builder(DSL)
var Dsl = &dsl{
	// Any value to StringGetter
	A: ToTextGetter,
	// Any objects of list ot TextList
	AL: ToTextList,
	// String value to StringGetter
	S: func(v string) TextGetter {
		return StringGetter(v)
	},
	// object of fmt.Stringer to StringGetter
	SER: NewStringerGetter,
	PF: TextGetterPrintf,
}

type dsl struct {
	A func(v interface{}) TextGetter
	AL func(v ...interface{}) TextList
	S func(v string) TextGetter
	SER func(v fmt.Stringer) *StringerGetter
	PF func(format string, a ...interface{}) TextGetter
}
