package check

import (
	"encoding/json"
	"fmt"
	"gopkg.in/check.v1"
	sjson "github.com/bitly/go-simplejson"
	"io"
)

var JsonEquals = &checkJsonEquals{}

type checkJsonEquals struct{}
func (c *checkJsonEquals) Info() *check.CheckerInfo {
	return &check.CheckerInfo {
		Name: "JsonEquals",
		Params: []string{ "obtained", "expected" },
	}
}
func (c *checkJsonEquals) Check(params []interface{}, names []string) (bool, string) {
	var obtained, expected =
		convertToSimpleJson(params[0]),
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
		[]interface{} { obtained.Interface(), expected.Interface() },
		names,
	)

	if !result {
		return false, fmt.Sprintf("Object JSON:[%v]. Expected JSON: [%v]", obtained, expected)
	}

	return true, ""
}

const (
	arrayType = 1
	mapType = 2
	valueType = 3
)

var typeToString = map[int]string {
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
