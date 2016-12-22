package json

import (
	json "github.com/bitly/go-simplejson"
)

type JsonExt struct {
	*json.Json
}

func ToJsonExt(json *json.Json) *JsonExt {
	return &JsonExt{ json }
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
	json, check := j.CheckGet(key)
	return &JsonExt{ json }, check
}

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
