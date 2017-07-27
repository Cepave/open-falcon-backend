package textbuilder

import (
	"fmt"
)

// Build-in double-quoted transformers
var QQ = map[string]Transformer{
	"()": BuildSurrounding(Dsl.S("("), Dsl.S(")")),
	"[]": BuildSurrounding(Dsl.S("["), Dsl.S("]")),
	"{}": BuildSurrounding(Dsl.S("{"), Dsl.S("}")),
	"<>": BuildSurrounding(Dsl.S("<"), Dsl.S(">")),
	"\"": BuildSameSurrounding(Dsl.S("\"")),
	"'":  BuildSameSurrounding(Dsl.S("'")),
}

// Common characters for distiller
var J = map[string]Distiller{
	",":  BuildJoin(Dsl.S(",")),
	", ": BuildJoin(Dsl.S(", ")),
}

// Short name of building blocks for text builder(DSL)
//
// 	A - Any value to StringGetter
// 	AL - Any list of objects to TextList
// 	S - String value to StringGetter
// 	SER - object of "fmt.Stringer" to StringerGetter
// 	PF - "fmt.Sprintf" talk to TextGetter
var Dsl = &dsl{
	A:  ToTextGetter,
	AL: ToTextList,
	S: func(v string) TextGetter {
		return StringGetter(v)
	},
	SER: NewStringerGetter,
	PF:  TextGetterPrintf,
}

type dsl struct {
	A   func(v interface{}) TextGetter
	AL  func(v ...interface{}) TextList
	S   func(v string) TextGetter
	SER func(v fmt.Stringer) *StringerGetter
	PF  func(format string, a ...interface{}) TextGetter
}
