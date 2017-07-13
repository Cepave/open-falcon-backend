package conform

// This library is a re-written version, which comes from:
// 	https://github.com/leebenson/conform
import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	or "github.com/Cepave/open-falcon-backend/common/reflect"
	"github.com/etgryphon/stringUp"
	e "github.com/juju/errors"
)

type ConformService struct {
	transformers map[string]StringTransformer

	cacheLock        *sync.Mutex
	transformerCache map[string][]StringTransformer
}

func NewConformService() *ConformService {
	return &ConformService{
		transformers: make(map[string]StringTransformer),

		cacheLock:        &sync.Mutex{},
		transformerCache: make(map[string][]StringTransformer),
	}
}

func (s *ConformService) RegisterStringTransformer(name string, transformer StringTransformer) {
	s.transformers[name] = transformer
}
func (s *ConformService) GetTransformer(name string) StringTransformer {
	transformer, ok := s.transformers[name]
	if !ok {
		transformer, ok = buildinTransformer[name]
		if !ok {
			panic(e.Errorf("Unknown name of transformer: %s", name))
		}
	}

	return transformer
}
func (s *ConformService) Conform(v interface{}) error {
	return s.conformAny(reflect.ValueOf(v))
}
func (s *ConformService) MustConform(v interface{}) {
	if err := s.Conform(v); err != nil {
		panic(e.ErrorStack(err))
	}
}

func (s *ConformService) conformAny(anyValue reflect.Value) error {
	return s.conformAnyWithTransformers(anyValue, nil)
}
func (s *ConformService) conformAnyWithTransformers(anyValue reflect.Value, transformers []StringTransformer) error {
	typeOfAnyValue := anyValue.Type()
	if typeOfAnyValue.Implements(_t_Conformer) {
		anyValue.Interface().(Conformer).ConformSelf(s)
		return nil
	}

	finalValueOfAny := or.FinalPointedValue(anyValue)
	finalTypeOfAnyValue := finalValueOfAny.Type()

	switch finalTypeOfAnyValue.Kind() {
	case reflect.Ptr:
		if finalValueOfAny.IsNil() {
			return nil
		}
	case reflect.Struct:
		for i := 0; i < finalValueOfAny.NumField(); i++ {
			field := finalTypeOfAnyValue.Field(i)
			transformersInTag := s.getTransformers(finalTypeOfAnyValue, field)
			if err := s.conformAnyWithTransformers(
				finalValueOfAny.Field(i), transformersInTag,
			); err != nil {
				return e.Annotatef(err, "Cannot conform field: %v", field)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < finalValueOfAny.Len(); i++ {
			elemValue := finalValueOfAny.Index(i)

			if err := s.conformAnyWithTransformers(elemValue, transformers); err != nil {
				return e.Annotatef(
					err, "Cannot conform elem[%d]. type: [%v]",
					i, typeOfAnyValue,
				)
			}
		}
	case reflect.String:
		return conformStringValue(anyValue, finalValueOfAny, transformers)
	default:
		if len(transformers) > 0 {
			return e.Errorf("Unsupported type: %v. Pointed type[%v]", typeOfAnyValue, finalTypeOfAnyValue)
		}
	}

	return nil
}
func (s *ConformService) getTransformers(structType reflect.Type, field reflect.StructField) []StringTransformer {
	uniqName := fmt.Sprintf("%s%s%s", structType.PkgPath(), structType.Name(), field.Name)
	if cachedTransformers, ok := s.transformerCache[uniqName]; ok {
		return cachedTransformers
	}

	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	/**
	 * Maybe some goroutin has built the cache before the acquiring of this lock
	 */
	if cachedTransformers, ok := s.transformerCache[uniqName]; ok {
		return cachedTransformers
	}
	// :~)

	tag := field.Tag.Get("conform")
	if tag == "" {
		return nil
	}

	transformers := make([]StringTransformer, 0)
	transformerNames := strings.Split(tag, ",")

	for _, name := range transformerNames {
		transformers = append(transformers, s.GetTransformer(name))
	}
	s.transformerCache[uniqName] = transformers

	return s.transformerCache[uniqName]
}

func conformStringValue(sourceValue reflect.Value, finalValue reflect.Value, transformers []StringTransformer) error {
	if len(transformers) == 0 {
		return nil
	}

	finalValue = or.FinalPointedValue(finalValue)

	if finalValue.Kind() != reflect.String {
		return e.Errorf("Value: %v is not string type: [%v]", finalValue, finalValue.Type())
	}

	stringValue := finalValue.Interface().(string)
	for _, t := range transformers {
		stringValue = t(stringValue)
	}

	if stringValue == nilString {
		if sourceValue.Kind() != reflect.Ptr {
			return e.Errorf("Cannot set xxxToNil to **NON-POINTER** type: %v", sourceValue.Type())
		}

		sourceValue.Set(reflect.Zero(sourceValue.Type()))
		return nil
	}

	finalValue.Set(reflect.ValueOf(stringValue))

	return nil
}

type StringTransformer func(string) string

// If a type implements this interface, the callback function gets called while conforming
type Conformer interface {
	ConformSelf(*ConformService)
}

var _t_Conformer = or.TypeOfInterface((*Conformer)(nil))

var _defaultService = NewConformService()

// Performs must conform with buildin transformers
func MustConform(v interface{}) {
	_defaultService.MustConform(v)
}

var patterns = map[string]*regexp.Regexp{
	"numbers":    regexp.MustCompile("[0-9]"),
	"nonNumbers": regexp.MustCompile("[^0-9]"),
	"alpha":      regexp.MustCompile("[\\pL]"),
	"nonAlpha":   regexp.MustCompile("[^\\pL]"),
	"name":       regexp.MustCompile("[\\p{L}]([\\p{L}|[:space:]|-]*[\\p{L}])*"),
}

const nilString = "<@!NilString!@>"

var buildinTransformer = map[string]StringTransformer{
	"trimToNil": func(input string) string {
		s := strings.TrimSpace(input)
		if s == "" {
			return nilString
		}

		return s
	},
	"trim":    strings.TrimSpace,
	"ltrim":   func(input string) string { return strings.TrimLeft(input, " ") },
	"rtrim":   func(input string) string { return strings.TrimRight(input, " ") },
	"lower":   strings.ToLower,
	"upper":   strings.ToUpper,
	"title":   strings.Title,
	"camel":   stringUp.CamelCase,
	"snake":   func(input string) string { return camelTo(stringUp.CamelCase(input), "_") },
	"slug":    func(input string) string { return camelTo(stringUp.CamelCase(input), "-") },
	"ucfirst": ucFirst,
	"name":    formatName,
	"email":   func(input string) string { return strings.ToLower(strings.TrimSpace(input)) },
	"num":     onlyNumbers,
	"!num":    stripNumbers,
	"alpha":   onlyAlpha,
	"!alpha":  stripAlpha,
}

func camelTo(s, sep string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" {
				words = append(words, initialism)

				i += len(initialism) - 1
				lastPos = i
				continue
			}

			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += sep
		}

		result += strings.ToLower(word)
	}

	return result
}

func ucFirst(s string) string {
	if s == "" {
		return s
	}
	toRune, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLower(toRune) {
		return s
	}
	buf := &bytes.Buffer{}
	buf.WriteRune(unicode.ToUpper(toRune))
	buf.WriteString(s[size:])
	return buf.String()
}

type x map[string]string

func formatName(s string) string {
	first := onlyOne(strings.ToLower(s), []x{{"[^\\pL-\\s]": ""}, {"\\s": " "}, {"-": "-"}})
	return strings.Title(patterns["name"].FindString(first))
}

func onlyNumbers(s string) string {
	return patterns["nonNumbers"].ReplaceAllLiteralString(s, "")
}

func stripNumbers(s string) string {
	return patterns["numbers"].ReplaceAllLiteralString(s, "")
}

func onlyAlpha(s string) string {
	return patterns["nonAlpha"].ReplaceAllLiteralString(s, "")
}

func stripAlpha(s string) string {
	return patterns["alpha"].ReplaceAllLiteralString(s, "")
}

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/3d26dc39376c307203d3a221bada26816b3073cf/lint.go#L482
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] {
			initialism = s[:i]
		}
	}
	return initialism
}

func onlyOne(s string, m []x) string {
	for _, v := range m {
		for f, r := range v {
			s = regexp.MustCompile(fmt.Sprintf("%s{2,}", f)).ReplaceAllLiteralString(s, r)
		}
	}
	return s
}
