package check

import (
	"encoding/json"
	"fmt"
	sjson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"io"
)

// Checks two varialbes which could be converted to "*go-simplejson.Json"
//
// Various types of variable are supported:
//
// 	string - JSON string
// 	[]byte - As JSON string
// 	io.Reader - A reader contains JSON
// 	*go-simplejson.Json
//
// If non of above types is matched, this checker uses "encoding/json.Marshal()" to marshal the object
// to JSON format and performs comparison.
//
// See "https://godoc.org/github.com/bitly/go-simplejson" for detail information of simplejson
var JsonEquals = &checkJsonEquals{}

type checkJsonEquals struct{}

func (c *checkJsonEquals) Info() *check.CheckerInfo {
	return &check.CheckerInfo{
		Name:   "JsonEquals",
		Params: []string{"obtained", "expected"},
	}
}
func (c *checkJsonEquals) Check(params []interface{}, names []string) (bool, string) {
	var obtained, expected = convertToSimpleJson(params[0]),
		convertToSimpleJson(params[1])

	obtainedType := getJsonType(obtained)
	expectedType := getJsonType(expected)

	if obtainedType != expectedType {
		return false, fmt.Sprintf(
			"Obtained type: [%s]. Expected type: [%s]",
			typeToString[obtainedType],
			typeToString[expectedType],
		)
	}

	result, _ := check.DeepEquals.Check(
		[]interface{}{obtained.Interface(), expected.Interface()},
		names,
	)

	if !result {
		return false, fmt.Sprintf("Object JSON:[%v]. Expected JSON: [%v]", obtained, expected)
	}

	return true, ""
}

const (
	arrayType = 1
	mapType   = 2
	valueType = 3
)

var typeToString = map[int]string{
	1: "JSON Array",
	2: "JSON Map",
	3: "JSON Value",
}

func getJsonType(jsonObject *sjson.Json) int {
	switch jsonObject.Interface().(type) {
	case []interface{}:
		return arrayType
	case map[string]interface{}:
		return mapType
	}

	return valueType
}

func convertToSimpleJson(v interface{}) *sjson.Json {
	for {
		switch jsonValue := v.(type) {
		case string:
			if v == "" {
				return nil
			}
			v = []byte(jsonValue)
		case []byte:
			jsonObj, objErr := sjson.NewJson(jsonValue)
			if objErr != nil {
				panic(fmt.Sprintf("Cannot convert string:[%s] to *simplejson.Json. %v", jsonValue, objErr))
			}

			return jsonObj
		case *sjson.Json:
			return jsonValue
		case io.Reader:
			jsonObj, objErr := sjson.NewFromReader(jsonValue)
			if objErr != nil {
				panic(fmt.Sprintf("Cannot convert Reader:[%T] to *simplejson.Json. %v", v, objErr))
			}
			return jsonObj
		default:
			jsonString, err := json.Marshal(jsonValue)
			if err != nil {
				panic(fmt.Sprintf("Cannot convert JSON:[%#v]. %v", jsonValue, err))
			}

			v = jsonString
		}
	}
}
