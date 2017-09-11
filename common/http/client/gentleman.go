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

// Converts "*gentleman.Response" to wrapper object of utility.
func ToGentlemanResp(resp *gt.Response) *GentlemanResponse {
	return (*GentlemanResponse)(resp)
}

// Wrapper object with providing additional functions to "*gentleman.Response"
type GentlemanResponse gt.Response

// Gets the json object(by go-simplejson) with error.
//
// Be careful: This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) GetJson() (*sjson.Json, error) {
	json := sjson.New()

	if err := r.BindJson(json); err != nil {
		return nil, errors.Annotate(err, "Building of JSON object is failed")
	}

	return json, nil
}

// Gets the json object(by go-simplejson) or panic if some error has occurred.
//
// Be careful: This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) MustGetJson() *sjson.Json {
	json, err := r.GetJson()
	if err != nil {
		panic(errors.Details(err))
	}

	return json
}

// Binds the body of response to input object(as JSON).
//
// Be careful: This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) BindJson(v interface{}) error {
	resp := (*gt.Response)(r)
	contentType := resp.Header.Get("Content-Type")

	return errors.Annotatef(
		resp.JSON(v),
		"Binding response body to JSON has error.\n\tStatus[%d] Content-Type: %s",
		resp.StatusCode, contentType,
	)
}

// Binds the body of response to input object(as JSON).
//
// This function would panic if some error has occurred.
//
// Be careful: This function would **close the response** because of the implementation of Gentleman library.
func (r *GentlemanResponse) MustBindJson(v interface{}) {
	if err := r.BindJson(v); err != nil {
		panic(errors.Details(err))
	}
}

// Gets the string content of body.
//
// The format of returned value:
//
// If the body is empty, gives:
//
// 	Response[<StatusCode>(<ContentType>)]
//
// If the body is viable, gives:
//
// 	Response[<StatusCode>](ContentType). Body: << <Shortened Body> >>
//
// The shortened body would be 128 characters of maximum by preserving prefix and suffix string of body.
//
// Be careful: This function would **close the response** because of the implementation of Gentleman library.
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

// Converts "*gentleman.Request" to "*GentlemanRequest"
func ToGentlemanReq(req *gt.Request) *GentlemanRequest {
	return (*GentlemanRequest)(req)
}

// This type is used with "*GentlemanRequest", which
// gets called for checking the expected response.
type RespMatcher func(*gt.Response) error

// Matcher for status code of HTTP
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

// Wrapper object with providing additional functions to "*gentleman.Request"
type GentlemanRequest gt.Request

// Sends request with expected status code
//
// If the status code is not as expected, the error would not be nil.
func (r *GentlemanRequest) SendAndStatusMatch(status int) (*gt.Response, error) {
	return r.SendAndMatch(StatusMatcher(status))
}

// Sends request with expected status code
//
// If the status code is not as expected, the function would be panic.
func (r *GentlemanRequest) SendAndStatusMustMatch(status int) *gt.Response {
	return r.SendAndMustMatch(StatusMatcher(status))
}

// Sends request with checking by implementation of matcher.
//
// If the status code is not as expected, the error would not be nil.
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

// Sends request with checking by implementation of matcher.
//
// If the status code is not as expected, the function would be panic.
func (r *GentlemanRequest) SendAndMustMatch(matcher RespMatcher) *gt.Response {
	resp, err := r.SendAndMatch(matcher)
	if err != nil {
		panic(errors.Details(err))
	}

	return resp
}

// Common configurations used for building of "*gentleman.Client" object.
type GentlemanConfig struct {
	RequestTimeout time.Duration
}

// Utility functions for building "*gentleman.Client" or "*gentleman.Request" with
// default configuration.
//
// 	Request Timeout: See "DEFAULT_TIMEOUT"
//
// NewDefaultClient()
// 	Constructs a client with default configuration.
//
// NewClientByConfig(*GentlemanConfig)
// 	Constructs a client with provided configuration.
//
// NewDefaultRequest()
// 	Constructs a request with default configuration.
//
// NewRequestByConfig(*GentlemanConfig)
// 	Constructs a request with provided configuration.
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
		RequestTimeout: DEFAULT_TIMEOUT,
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
