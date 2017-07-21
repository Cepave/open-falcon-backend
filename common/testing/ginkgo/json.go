package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

func MatchJson(v interface{}) GomegaMatcher {
	return &matchJsonImpl{
		MatchJSON(loadJsonString(v)),
	}
}

type matchJsonImpl struct {
	gomegaMatcher GomegaMatcher
}

func (m *matchJsonImpl) Match(actual interface{}) (success bool, err error) {
	jsonContent := loadJsonString(actual)
	return m.gomegaMatcher.Match(jsonContent)
}
func (m *matchJsonImpl) FailureMessage(actual interface{}) (message string) {
	jsonContent := loadJsonString(actual)
	return m.gomegaMatcher.FailureMessage(jsonContent)
}
func (m *matchJsonImpl) NegatedFailureMessage(actual interface{}) (message string) {
	jsonContent := loadJsonString(actual)
	return m.gomegaMatcher.FailureMessage(jsonContent)
}

func loadJsonString(v interface{}) interface{} {
	if jsonMarshaler, ok := v.(json.Marshaler); ok {
		jsonContent, err := jsonMarshaler.MarshalJSON()
		if err != nil {
			panic(err)
		}

		return jsonContent
	}

	return v
}
