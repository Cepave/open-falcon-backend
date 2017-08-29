package client

import (
	"fmt"
	"time"

	sjson "github.com/bitly/go-simplejson"
	"github.com/h2non/gentleman/plugins/timeout"
	"github.com/juju/errors"
	gt "gopkg.in/h2non/gentleman.v2"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

func ToGentlemanResp(resp *gt.Response) *GentlemanResponse {
	return (*GentlemanResponse)(resp)
}

type GentlemanResponse gt.Response

// This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) GetJson() (*sjson.Json, error) {
	json := sjson.New()

	if err := r.BindJson(json); err != nil {
		return nil, errors.Annotate(err, "Building of JSON object is failed")
	}

	return json, nil
}

// This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) MustGetJson() *sjson.Json {
	json, err := r.GetJson()
	if err != nil {
		panic(errors.Details(err))
	}

	return json
}

// This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) BindJson(v interface{}) error {
	resp := (*gt.Response)(r)
	contentType := resp.Header.Get("Content-Type")

	return errors.Annotatef(
		resp.JSON(v),
		"Binding response body to JSON has error.\n\tStatus[%d] Content-Type: %s",
		resp.StatusCode, contentType,
	)
}

// This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) MustBindJson(v interface{}) {
	if err := r.BindJson(v); err != nil {
		panic(errors.Details(err))
	}
}

// !!Be Caution!! This function would change the body to "closed".
func (r *GentlemanResponse) ToDetailString() string {
	resp := (*gt.Response)(r)

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "!UNKNOWN!"
	}

	bodyString := resp.String()
	defer resp.ClearInternalBuffer()

	if bodyString == "" {
		return fmt.Sprintf("Response[%d](%s)", resp.StatusCode, contentType)
	}

	bodyString = utils.ShortenStringToSize(bodyString, " ... ", 128)
	return fmt.Sprintf("Response[%d](%s). Body: << %s >>", resp.StatusCode, contentType, bodyString)
}

func ToGentlemanReq(req *gt.Request) *GentlemanRequest {
	return (*GentlemanRequest)(req)
}

type RespMatcher func(*gt.Response) error

func StatusMatcher(status int) RespMatcher {
	return func(resp *gt.Response) error {
		if resp.StatusCode == 0 {
			return errors.Errorf("HTTP status(0, may be timeout) of response is not matched to [%d].", status)
		}

		if resp.StatusCode != status {
			return errors.Errorf("HTTP status[%d] of response is not matched to [%d].", resp.StatusCode, status)
		}

		return nil
	}
}

type GentlemanRequest gt.Request

func (r *GentlemanRequest) SendAndStatusMatch(status int) (*gt.Response, error) {
	return r.SendAndMatch(StatusMatcher(status))
}
func (r *GentlemanRequest) SendAndStatusMustMatch(status int) *gt.Response {
	return r.SendAndMustMatch(StatusMatcher(status))
}
func (r *GentlemanRequest) SendAndMatch(matcher RespMatcher) (*gt.Response, error) {
	request := (*gt.Request)(r)

	resp, err := request.Send()

	if err != nil {
		if resp != nil {
			defer resp.Close()
		}
		return nil, errors.Annotate(err, "Send() of gentleman request has error")
	}

	if matcherErr := matcher(resp); matcherErr != nil {
		defer resp.Close()

		return nil, errors.Errorf("HTTP response has error: %v", errors.Details(matcherErr))
	}

	return resp, nil
}
func (r *GentlemanRequest) SendAndMustMatch(matcher RespMatcher) *gt.Response {
	resp, err := r.SendAndMatch(matcher)
	if err != nil {
		panic(errors.Details(err))
	}

	return resp
}

// Common configurations used in
type GentlemanConfig struct {
	RequestTimeout time.Duration
}

var CommonGentleman = &struct {
	NewDefaultClient   func() *gt.Client
	NewClientByConfig  func(config *GentlemanConfig) *gt.Client
	NewDefaultRequest  func() *gt.Request
	NewRequestByConfig func(config *GentlemanConfig) *gt.Request
}{
	NewDefaultClient:   gtNewDefaultClient,
	NewClientByConfig:  gtNewClientByConfig,
	NewDefaultRequest:  gtNewDefaultRequest,
	NewRequestByConfig: gtNewRequestByConfig,
}

func gtNewDefaultClient() *gt.Client {
	return gtNewClientByConfig(&GentlemanConfig{
		RequestTimeout: time.Second * 10,
	})
}

func gtNewClientByConfig(config *GentlemanConfig) *gt.Client {
	return gt.New().Use(timeout.Request(config.RequestTimeout))
}

func gtNewDefaultRequest() *gt.Request {
	return gtNewDefaultClient().Request()
}

func gtNewRequestByConfig(config *GentlemanConfig) *gt.Request {
	return gtNewClientByConfig(config).Request()
}
