package textbuilder

import (
	"unicode/utf8"
	"fmt"
	"reflect"
)

const EmptyGetter = StringGetter("")

type TextGetter interface {
	fmt.Stringer
	Post
}

// Used to transform TextGetter
type Transformer func(TextGetter) TextGetter

// Used to generates TextList
type Breeder func(TextGetter) TextList

// Converts a list of text to TextGetter
type Distiller func(TextList) TextGetter

// Gets posting processor of a TextGetter
type Post interface {
	Post() PostProcessor
}

// Defines the operations of post processor
type PostProcessor interface {
	TextGetter
	Transform(t Transformer) PostProcessor
	Breed(b Breeder) TextList
	Prefix(prefix TextGetter) PostProcessor
	Suffix(suffix TextGetter) PostProcessor
	Surrounding(prefix TextGetter, suffix TextGetter) PostProcessor
	Repeat(times int) TextList
	RepeatByLen(lenObject interface{}) TextList
}

// Initialize a new instance of post processor with default operations
func NewPost(content TextGetter) *DefaultPost {
	return &DefaultPost{ content }
}
// Implements default post prcessor
type DefaultPost struct {
	content TextGetter
}
func (p *DefaultPost) Transform(t Transformer) PostProcessor {
	p.content = t(p.content)
	return p
}
func (p *DefaultPost) Breed(b Breeder) TextList {
	return b(p.content)
}
func (p *DefaultPost) Prefix(prefix TextGetter) PostProcessor {
	p.content = Prefix(prefix, p.content)
	return p
}
func (p *DefaultPost) Suffix(suffix TextGetter) PostProcessor {
	p.content = Suffix(p.content, suffix)
	return p
}
func (p *DefaultPost) Surrounding(prefix TextGetter, suffix TextGetter) PostProcessor {
	p.content = Surrounding(prefix, p.content, suffix)
	return p
}
func (p *DefaultPost) Repeat(times int) TextList {
	return Repeat(p.content, times)
}
func (p *DefaultPost) RepeatByLen(lenObject interface{}) TextList {
	return RepeatByLen(p.content, lenObject)
}
func (p *DefaultPost) Viable(v bool) PostProcessor {
	if !v {
		p.content = EmptyGetter
	}
	return p
}
func (p *DefaultPost) Post() PostProcessor {
	return p
}
func (p *DefaultPost) String() string {
	return p.content.String()
}

// Implements the text getter with string value
type StringGetter string
func (t StringGetter) String() string {
	return string(t)
}
func (t StringGetter) Post() PostProcessor {
	return NewPost(t)
}

// Converts fmt.Stringer interface to TextGetter
func NewStringerGetter(v fmt.Stringer) *StringerGetter {
	return &StringerGetter{ v }
}
type StringerGetter struct {
	stringer fmt.Stringer
}
func (s *StringerGetter) String() string {
	return s.stringer.String()
}
func (s *StringerGetter) Post() PostProcessor {
	return NewPost(s)
}

// Used to get len of an object
//
// This interface is usually used with RepeatByLen().
type ObjectLen interface {
	Len() int
}

type TextList interface {
	ListPost
	ObjectLen
	Get(int) TextGetter
}

// Gets posting processor of a TextList
type ListPost interface {
	Post() ListPostProcessor
}

// Defines operations for a TextList
type ListPostProcessor interface {
	Distill(Distiller) TextGetter
	Join(seperator TextGetter) TextGetter
}

// Initialzie an instance of DefaultListPost
func NewListPost(list TextList) *DefaultListPost {
	return &DefaultListPost { list }
}

// Implements default post prcessor for a list
type DefaultListPost struct {
	list TextList
}
func (l *DefaultListPost) Join(seperator TextGetter) TextGetter {
	return JoinTextList(seperator, l.list)
}
func (l *DefaultListPost) Distill(d Distiller) TextGetter {
	return d(l.list)
}

type TextGetters []TextGetter
func Getters(getters ...TextGetter) TextGetters {
	return TextGetters(getters)
}
func (t TextGetters) Get(index int) TextGetter {
	return t[index]
}
func (t TextGetters) Len() int {
	return len(t)
}
func (t TextGetters) Post() ListPostProcessor {
	return NewListPost(t)
}

// Converts any value to TextGetter
//
// If the value is text getter, this function return it natively.
//
// If the value is string, this function cast it to StringGetter.
//
// Otherwise, use fmt.Sprintf("%v") to retrieve the string representation of input value.
func ToTextGetter(v interface{}) TextGetter {
	switch castedValue := v.(type) {
	case TextGetter:
		return castedValue
	case fmt.Stringer:
		return NewStringerGetter(v.(fmt.Stringer))
	case string:
		return StringGetter(castedValue)
	}

	return TextGetterPrintf("%v", v)
}

// Converts multiple values to TextList, for the convertion of element,  see ToTextGetter
func ToTextList(anyObjects ...interface{}) TextList {
	getters := make([]TextGetter, len(anyObjects))

	for i, v := range anyObjects {
		getters[i] = ToTextGetter(v)
	}

	return TextGetters(getters)
}

func TextGetterPrintf(format string, a ...interface{}) TextGetter {
	return &formatterImpl{format, a}
}

// Builds viable getter if the value is viable
//
// The value would be evaluated eagerly.
//
// 	For string - must be non-empty
// 	For TextGetter - the result of content must be non empty
// 	For array, slice, map, chan - the len(array) > 0
//
// 	Otherwise - value.IsNil() should be false
func IsViable(value interface{}) bool {
	switch textValue := value.(type) {
	case string:
		return textValue != ""
	case fmt.Stringer:
		return IsViable(NewStringerGetter(textValue).String())
	case TextGetter:
		return IsViable(textValue.String())
	}

	reflectValue := reflect.ValueOf(value)

	switch reflectValue.Type().Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return reflectValue.Len() > 0
	}

	return !reflectValue.IsNil()
}

// Prefixing the content(if the content viable)
func Prefix(prefix TextGetter, content TextGetter) TextGetter {
	return &prefixImpl { prefix, content }
}

// Suffixing the content(if the content is viable)
func Suffix(content TextGetter, suffix TextGetter) TextGetter {
	return &suffixImpl { content, suffix }
}

// Surrounding the content(if the content is viable)
func Surrounding(prefix TextGetter, content TextGetter, suffix TextGetter) TextGetter {
	return &surroundingImpl { prefix, content, suffix }
}

// Surrounding the content(if the content is viable)
func SurroundingSame(s TextGetter, content TextGetter) TextGetter {
	return &surroundingImpl { s, content, s }
}

// Joining the viable element of getters
func Join(seperator TextGetter, getters ...TextGetter) TextGetter {
	return JoinTextList(seperator, TextGetters(getters))
}

// Joining the viable element of TextList
func JoinTextList(seperator TextGetter, textList TextList) TextGetter {
	return &joinImpl { seperator, textList }
}

// Repeating the viable element of TextList
func Repeat(text TextGetter, times int) TextList {
	list := make(TextGetters, times)

	for i := 0; i < times; i++ {
		list[i] = text
	}

	return list
}

// Repeats the len of object:
//
// For object len: use Len() function
// For String: use utf8.RuneCountInString(<string>) function
// For Array, Chan, Map, or Slice: use reflect.Value.Len() function
func RepeatByLen(text TextGetter, lenObject interface{}) TextList {
	var repeatTimes int

	switch v := lenObject.(type) {
	case ObjectLen:
		repeatTimes = v.Len()
	case string:
		repeatTimes = utf8.RuneCountInString(v)
	default:
		value := reflect.ValueOf(lenObject)
		switch value.Kind() {
			case reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
				repeatTimes = value.Len()
			default:
				panic(fmt.Sprintf("Cannot figure out the \"len\" of type[%T].", lenObject))
		}
	}

	return Repeat(text, repeatTimes)
}

func RepeatAndJoin(text TextGetter, seperator TextGetter, times int) TextGetter {
	return JoinTextList(seperator, Repeat(text, times));
}
func RepeatAndJoinByLen(text TextGetter, seperator TextGetter, lenObject interface{}) TextGetter {
	return JoinTextList(seperator, RepeatByLen(text, lenObject));
}

type formatterImpl struct {
	formatter string
	args []interface{}
}
func (f *formatterImpl) String() string {
	return fmt.Sprintf(f.formatter, f.args...)
}
func (f *formatterImpl) Post() PostProcessor {
	return NewPost(f)
}
