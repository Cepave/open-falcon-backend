package ginkgo

import (
	ohttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	"io"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests MatchHttpStatus", func() {
	resp := &http.Response{StatusCode: 200}

	DescribeTable("Matching result is true",
		func(actual interface{}) {
			Expect(actual).To(MatchHttpStatus(200))
		},
		Entry("By *http.Response", resp),
		Entry("By *testing/http.ResponseResult",
			&ohttp.ResponseResult{Response: resp},
		),
	)

	DescribeTable("Matching result is false",
		func(actual interface{}) {
			Expect(actual).ToNot(MatchHttpStatus(400))
		},
		Entry("By *http.Response", resp),
		Entry("By *testing/http.ResponseResult",
			&ohttp.ResponseResult{Response: resp},
		),
	)
})

var _ = Describe("Tests MatchHttpBodyAsJson", func() {
	sampleJson := `{ "name": "joe", "age": 33 }`

	newResp := func() *http.Response {
		return &http.Response{Body: &stringCloser{strings.NewReader(sampleJson)}}
	}

	DescribeTable("Matching result is true",
		func(actual interface{}) {
			Expect(actual).To(MatchHttpBodyAsJson(sampleJson))
		},
		Entry("By *http.Response", newResp()),
		Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(newResp())),
	)

	DescribeTable("Matching result is false",
		func(actual interface{}) {
			Expect(actual).ToNot(MatchHttpBodyAsJson(`{ "name": "joe", "age": 34 }`))
		},
		Entry("By *http.Response", newResp()),
		Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(newResp())),
	)
})

type stringCloser struct {
	io.Reader
}

func (s *stringCloser) Close() error {
	return nil
}
