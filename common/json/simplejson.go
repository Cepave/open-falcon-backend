package json

import (
	gjson "encoding/json"
	sjson "github.com/bitly/go-simplejson"
)

type SimpleJsonMarshaler interface {
	MarshalSimpleJSON() (*sjson.Json, error)
}
type SimpleJsonUnmarshaler interface {
	UnmarshalSimpleJSON(*sjson.Json) error
}

type JsonExt struct {
	*sjson.Json
}

func UnmarshalToJson(v interface{}) *sjson.Json {
	var jsonString []byte

	switch typedJson := v.(type) {
	case []byte:
		jsonString = typedJson
	case string:
		jsonString = []byte(typedJson)
	case gjson.Marshaler:
		marshalerJson, err := typedJson.MarshalJSON()
		if err != nil {
			panic(err)
		}

		jsonString = marshalerJson
	default:
		anyJson, err := gjson.Marshal(typedJson)
		if err != nil {
			panic(err)
		}

		jsonString = anyJson
	}

	jsonObject := sjson.New()

	if len(jsonString) == 0 {
		return nil
	}

	err := jsonObject.UnmarshalJSON(jsonString)
	if err != nil {
		panic(err)
	}

	return jsonObject
}

func UnmarshalToJsonExt(v interface{}) *JsonExt {
	jsonObject := UnmarshalToJson(v)
	if jsonObject == nil {
		return nil
	}

	return ToJsonExt(jsonObject)
}

func ToJsonExt(sjson *sjson.Json) *JsonExt {
	return &JsonExt{ sjson }
}

func (j *JsonExt) MustInt8() int8 {
	return int8(j.MustInt64())
}
func (j *JsonExt) MustInt16() int16 {
	return int16(j.MustInt64())
}
func (j *JsonExt) MustInt32() int32 {
	return int32(j.MustInt64())
}

func (j *JsonExt) MustUint8() uint8 {
	return uint8(j.MustUint64())
}
func (j *JsonExt) MustUint16() uint16 {
	return uint16(j.MustUint64())
}
func (j *JsonExt) MustUint32() uint32 {
	return uint32(j.MustUint64())
}

func (j *JsonExt) GetExt(key string) *JsonExt {
	return &JsonExt{ j.Get(key) }
}
func (j *JsonExt) GetIndexExt(index int) *JsonExt {
	return &JsonExt{ j.GetIndex(index) }
}
func (j *JsonExt) GetPathExt(branch ...string) *JsonExt {
	return &JsonExt{ j.GetPath(branch...) }
}
func (j *JsonExt) CheckGetExt(key string) (*JsonExt, bool) {
	sjson, check := j.CheckGet(key)
	return &JsonExt{ sjson }, check
}

// Gets the JSON
//
// This method would panic if the JSON cannot be marshalled
func MarshalJSON(anyObject interface{}) string {
	jsonString, err := gjson.Marshal(anyObject)
	if err != nil {
		panic(err)
	}

	return string(jsonString)
}

// Gets the JSON(pretty)
//
// This method would panic if the JSON cannot be marshalled
func MarshalPrettyJSON(anyObject interface{}) string {
	jsonString := MarshalJSON(anyObject)

	jsonObject, err := sjson.NewJson([]byte(jsonString))
	if err != nil {
		panic(err)
	}

	prettyJson, jsonError := jsonObject.EncodePretty()
	if jsonError != nil {
		panic(jsonError)
	}

	return string(prettyJson)
}
