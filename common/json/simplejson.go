package json

import (
	gjson "encoding/json"
	sjson "github.com/bitly/go-simplejson"
)

type JsonExt struct {
	*sjson.Json
}

func MarshalToJsonExt(v interface{}) *JsonExt {
	switch jsonObj := v.(type) {
	case *sjson.Json:
		return ToJsonExt(jsonObj)
	}

	jsonString, err := gjson.Marshal(v)
	if err != nil {
		panic(err)
	}

	sjsonObject, err := sjson.NewJson([]byte(jsonString))
	if err != nil {
		panic(err)
	}

	return ToJsonExt(sjsonObject)
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
func MarshalJSON(jsonContent *sjson.Json) string {
	return MarshalGoJson(jsonContent)
}

func MarshalAny(v interface{}) string {
	switch jsonObj := v.(type) {
	case *sjson.Json:
		return MarshalJSON(jsonObj)
	case gjson.Marshaler:
		return MarshalGoJson(jsonObj)
	}

	jsonString, err := gjson.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(jsonString)
}

func MarshalGoJson(jsonMarshaler gjson.Marshaler) string {
	jsonByte, err := jsonMarshaler.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(jsonByte)
}

// Gets the JSON(pretty)
//
// This method would panic if the JSON cannot be marshalled
func MarshalPrettyJSON(jsonContent *sjson.Json) string {
	jsonByte, err := jsonContent.EncodePretty()
	if err != nil {
		panic(err)
	}

	return string(jsonByte)
}
