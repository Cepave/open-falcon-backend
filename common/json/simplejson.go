package json

import (
	"io"

	gjson "encoding/json"

	sjson "github.com/bitly/go-simplejson"

	"github.com/juju/errors"
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
	switch typedJson := v.(type) {
	case *sjson.Json:
		return typedJson
	case SimpleJsonMarshaler:
		newJson, err := typedJson.MarshalSimpleJSON()
		if err != nil {
			err = errors.Annotate(err, "SimpleJsonMarshaler.MarshalSimpleJSON() has error")
			panic(errors.Details(err))
		}
		return newJson
	case gjson.Marshaler:
		marshalerJson, err := typedJson.MarshalJSON()
		if err != nil {
			err = errors.Annotate(err, "encoding/json.MarshalJSON() has error")
			panic(errors.Details(err))
		}

		return UnmarshalToJson(marshalerJson)
	case string:
		return UnmarshalToJson([]byte(typedJson))
	case []byte:
		if len(typedJson) == 0 {
			return sjson.New()
		}
		jsonObject, err := sjson.NewJson(typedJson)
		if err != nil {
			err = errors.Annotate(err, "go-simplejson.NewJson([]byte) has error")
			panic(errors.Details(err))
		}
		return jsonObject
	case io.Reader:
		jsonObject, err := sjson.NewFromReader(typedJson)
		if err != nil {
			err = errors.Annotate(err, "go-simplejson.NewFromReader(io.Reader) has error")
			panic(errors.Details(err))
		}
		return jsonObject
	}

	anyJson, err := gjson.Marshal(v)
	if err != nil {
		err = errors.Annotate(err, "encoding/json.Marshal() has error")
		panic(errors.Details(err))
	}

	return UnmarshalToJson(anyJson)
}

func UnmarshalToJsonExt(v interface{}) *JsonExt {
	jsonObject := UnmarshalToJson(v)
	if jsonObject == nil {
		return nil
	}

	return ToJsonExt(jsonObject)
}

func ToJsonExt(sjson *sjson.Json) *JsonExt {
	return &JsonExt{sjson}
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

func (j *JsonExt) IsNil() bool {
	return j.Interface() == nil
}

func (j *JsonExt) MustStringPtr() *string {
	if j.IsNil() {
		return nil
	}

	stringValue := j.MustString()
	return &stringValue
}

func (j *JsonExt) GetExt(key string) *JsonExt {
	return &JsonExt{j.Get(key)}
}
func (j *JsonExt) GetIndexExt(index int) *JsonExt {
	return &JsonExt{j.GetIndex(index)}
}
func (j *JsonExt) GetPathExt(branch ...string) *JsonExt {
	return &JsonExt{j.GetPath(branch...)}
}
func (j *JsonExt) CheckGetExt(key string) (*JsonExt, bool) {
	sjson, check := j.CheckGet(key)
	return &JsonExt{sjson}, check
}

// Gets the JSON
//
// This method would panic if the JSON cannot be marshalled
func MarshalJSON(anyObject interface{}) string {
	jsonString, err := gjson.Marshal(anyObject)
	if err != nil {
		err = errors.Annotate(err, "encoding/json.Marshal() has error")
		panic(errors.Details(err))
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
		err = errors.Annotate(err, "simplejson.NewJson([]byte) has error")
		panic(errors.Details(err))
	}

	prettyJson, jsonError := jsonObject.EncodePretty()
	if jsonError != nil {
		err = errors.Annotate(err, "simplejson.JSON.EncodePretty() has error")
		panic(errors.Details(err))
	}

	return string(prettyJson)
}
