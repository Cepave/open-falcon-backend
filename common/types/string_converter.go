package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Defines the parsing of a string value(with locale) to desired value
type Parser interface {
	Scan(stringValue string, locale *time.Location) interface{}
}

// Defines the string-representation of a value to string value(with locale)
type Printer interface {
	Print(object interface{}, locale *time.Location) string
}

// Combines Parser and Printer interface
type Formatter interface {
	Parser
	Printer
}

// Sub-type of ConverterRegistry, with registering of formatter.
type FormatterRegistry interface {
	ConverterRegistry

	// Adds formatter with type of object
	AddFormatter(objectType reflect.Type, formatter Formatter)
	// Adds formatter by type of struct's field
	AddFormatterForField(field reflect.StructField, formatter Formatter)
}

// Sub-type of ConversionService, which provides methods of "Print()" and "Scan()"...
type FormattedConversionService interface {
	ConversionService

	// Generates the string-representation of an object
	Print(sourceObject interface{}) string
	// Generates the string-representation of an object(with locale)
	PrintWithLocale(sourceObject interface{}, locale *time.Location) string

	// Converts the string value to target value
	Scan(stringValue string) interface{}
	// Converts the string value to target value(with locale)
	ScanWithLocale(stringValue string, locale *time.Location) interface{}
}

func convertStringToFloat(object interface{}, targetType reflect.Type) interface{} {
	strValue := object.(string)
	if strValue == "" {
		strValue = "0"
	}

	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		panic(err)
	}

	switch targetType.Kind() {
	case reflect.Float32:
		return float32(floatValue)
	case reflect.Float64:
		return floatValue
	}

	panic(fmt.Sprintf("Unknown kind[%v] for float value", targetType.Kind()))
}
func convertStringToInt(object interface{}, targetType reflect.Type) interface{} {
	strValue := object.(string)
	if strValue == "" {
		strValue = "0"
	}

	intValue, err := strconv.ParseInt(strValue, 10, targetType.Bits())
	if err != nil {
		panic(err)
	}

	switch targetType.Kind() {
	case reflect.Int:
		return int(intValue)
	case reflect.Int8:
		return int8(intValue)
	case reflect.Int16:
		return int16(intValue)
	case reflect.Int32:
		return int32(intValue)
	case reflect.Int64:
		return intValue
	}

	panic(fmt.Sprintf("Unknown kind[%v] for int value", targetType.Kind()))
}
func convertStringToUint(object interface{}, targetType reflect.Type) interface{} {
	strValue := object.(string)
	if strValue == "" {
		strValue = "0"
	}

	uintValue, err := strconv.ParseUint(strValue, 10, targetType.Bits())
	if err != nil {
		panic(err)
	}

	switch targetType.Kind() {
	case reflect.Uint:
		return uint(uintValue)
	case reflect.Uint8:
		return uint8(uintValue)
	case reflect.Uint16:
		return uint16(uintValue)
	case reflect.Uint32:
		return uint32(uintValue)
	case reflect.Uint64:
		return uintValue
	}

	panic(fmt.Sprintf("Unknown kind[%v] for int value", targetType.Kind()))
}
func convertStringToBool(object interface{}, targetType reflect.Type) interface{} {
	s := object.(string)

	switch strings.ToLower(s) {
	case "y", "yes", "t", "true":
		return true
	default:
		intValue, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return intValue != 0
		}
	}

	return false
}
