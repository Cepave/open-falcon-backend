package json

import (
	json "github.com/bitly/go-simplejson"
)

// Gets the JSON
//
// This method would panic if the JSON cannot be marshalled
func MarshalJSON(jsonContent *json.Json) string {
	jsonByte, err := jsonContent.Encode()
	if err != nil {
		panic(err)
	}

	return string(jsonByte)
}

// Gets the JSON(pretty)
//
// This method would panic if the JSON cannot be marshalled
func MarshalPrettyJSON(jsonContent *json.Json) string {
	jsonByte, err := jsonContent.EncodePretty()
	if err != nil {
		panic(err)
	}

	return string(jsonByte)
}
