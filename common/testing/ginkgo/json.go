package ginkgo

import (
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	sjson "github.com/bitly/go-simplejson"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

func MatchJson(v interface{}) GomegaMatcher {
	expectedJsonString, _ := ojson.UnmarshalToJson(v).MarshalJSON()

	return &matchJsonImpl{
		matcher: MatchJSON(expectedJsonString),
	}
}

type matchJsonImpl struct {
	matcher    GomegaMatcher
	actualJson *sjson.Json
}

func (m *matchJsonImpl) Match(actual interface{}) (success bool, err error) {
	m.actualJson = ojson.UnmarshalToJson(actual)

	actualJsonString, _ := m.actualJson.MarshalJSON()

	return m.matcher.Match(actualJsonString)
}
func (m *matchJsonImpl) FailureMessage(actual interface{}) (message string) {
	actualJsonString, _ := m.actualJson.MarshalJSON()
	return m.matcher.FailureMessage(actualJsonString)
}
func (m *matchJsonImpl) NegatedFailureMessage(actual interface{}) (message string) {
	actualJsonString, _ := m.actualJson.MarshalJSON()
	return m.matcher.NegatedFailureMessage(actualJsonString)
}
