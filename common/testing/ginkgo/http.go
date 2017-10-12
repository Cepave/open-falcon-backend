package ginkgo

import (
	"fmt"
	"net/http"

	ohttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	gt "gopkg.in/h2non/gentleman.v2"

	. "github.com/onsi/gomega/types"
)

// Matches the status of HTTP response, the type of tested value could be:
// 	*http.Response
// 	*testing/http.ResponseResult
func MatchHttpStatus(status int) GomegaMatcher {
	return &httpStatusMatcher{status}
}

// Matches the JSON body of HTTP response, the type of tested value could be:
// 	*http.Response
// 	*testing/http.ResponseResult
func MatchHttpBodyAsJson(json interface{}) GomegaMatcher {
	return &jsonBodyMatcher{
		expectedBody: json,
		matcher:      MatchJson(json),
	}
}

type jsonBodyMatcher struct {
	expectedBody interface{}
	matcher      GomegaMatcher
}

func (m *jsonBodyMatcher) Match(actual interface{}) (success bool, err error) {
	respResult := getResponseResult(actual)
	if respResult == nil {
		return false, buildRespError(actual)
	}

	return m.matcher.Match(respResult.GetBodyAsJson())
}
func (m *jsonBodyMatcher) FailureMessage(actual interface{}) (message string) {
	respResult := getResponseResult(actual)

	return fmt.Sprintf("[HTTP] Expected JSON body is not match. %s",
		m.matcher.FailureMessage(respResult.GetBodyAsJson()),
	)
}
func (m *jsonBodyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	respResult := getResponseResult(actual)

	return fmt.Sprintf("[HTTP] Not expected match body of JSON. %s",
		m.matcher.NegatedFailureMessage(respResult.GetBodyAsJson()),
	)
}

type httpStatusMatcher struct {
	expectedStatus int
}

func (m *httpStatusMatcher) Match(actual interface{}) (success bool, err error) {
	resp := getResponse(actual)
	if resp == nil {
		return false, buildRespError(actual)
	}

	return resp.StatusCode == m.expectedStatus, nil
}
func (m *httpStatusMatcher) FailureMessage(actual interface{}) (message string) {
	resp := getResponse(actual)

	return fmt.Sprintf(
		"[HTTP] Expected status: [%d]. Got [%d](%s)",
		m.expectedStatus, resp.StatusCode, resp.Status,
	)
}
func (m *httpStatusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	resp := getResponse(actual)

	return fmt.Sprintf(
		"[HTTP] Not expected status: [%d](%s). But got it.",
		resp.StatusCode, resp.StatusCode,
	)
}

func buildRespError(actual interface{}) error {
	return fmt.Errorf(`Actual value need to be type of "*http.Response" or "testing/http.ResponseResult". Got: [%T]`, actual)
}
func getResponseResult(actual interface{}) *ohttp.ResponseResult {
	switch v := actual.(type) {
	case *ohttp.ResponseResult:
		return v
	case *http.Response:
		return ohttp.NewResponseResultByResponse(v)
	}

	return nil
}
func getResponse(actual interface{}) *http.Response {
	switch v := actual.(type) {
	case *ohttp.ResponseResult:
		return v.Response
	case *gt.Response:
		return v.RawResponse
	case *http.Response:
		return v
	}

	panic(fmt.Sprintf("Unsupported type of response: %T", actual))
}
