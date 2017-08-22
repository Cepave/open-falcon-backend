package client

import (
	"time"

	sjson "github.com/bitly/go-simplejson"
	"github.com/h2non/gentleman/plugins/timeout"
	"github.com/juju/errors"
	gt "gopkg.in/h2non/gentleman.v2"
)

func ToGentlemanResp(resp *gt.Response) *GentlemanResponse {
	return (*GentlemanResponse)(resp)
}

type GentlemanResponse gt.Response

func (r *GentlemanResponse) GetJson() (*sjson.Json, error) {
	json := sjson.New()

	err := (*gt.Response)(r).JSON(json)
	if err != nil {
		return nil, errors.Annotate(err, "Binding JSON of response has error")
	}

	return json, nil
}
func (r *GentlemanResponse) MustGetJson() *sjson.Json {
	json, err := r.GetJson()
	if err != nil {
		panic(err)
	}

	return json
}

func ToGentlemanReq(req *gt.Request) *GentlemanRequest {
	return (*GentlemanRequest)(req)
}

type GentlemanRequest gt.Request

func (r *GentlemanRequest) SendAndStatusMatch(status int) (*gt.Response, error) {
	request := (*gt.Request)(r)

	resp, err := request.Send()

	if err != nil {
		return nil, errors.Annotate(err, "Send() of gentleman request has error")
	}

	if resp.StatusCode != status {
		return nil, errors.Annotatef(err, "HTTP status of response is not match[%d]. Current: [%d]", status, resp.StatusCode)
	}

	return resp, nil
}
func (r *GentlemanRequest) SendAndStatusMustMatch(status int) *gt.Response {
	resp, err := r.SendAndStatusMatch(status)
	if err != nil {
		panic(err)
	}

	return resp
}

// Common configurations used in
type GentlemanConfig struct {
	RequestTimeout time.Duration
}

// Namespace of common functions for h2non/gentleman library
type GentlemanFuncs interface {
	// Default values:
	// 	Timeout(whole request) - 10 seconds
	NewDefaultClient() *gt.Client
	NewClientByConfig(config *GentlemanConfig) *gt.Client

	// Constructs a request by default values
	//
	// See NewDefaultClient() for default values of configuration.
	NewDefaultRequest() *gt.Request
	NewRequestByConfig(config *GentlemanConfig) *gt.Request
}

var CommonGentleman GentlemanFuncs = &gentlemanImpl{}

// This type is used for functions-aggregation only,
// MUST NOT HAS STATUS
type gentlemanImpl struct{}

func (g *gentlemanImpl) NewDefaultClient() *gt.Client {
	return g.NewClientByConfig(&GentlemanConfig{
		RequestTimeout: time.Second * 10,
	})
}

func (g *gentlemanImpl) NewClientByConfig(config *GentlemanConfig) *gt.Client {
	return gt.New().Use(timeout.Request(config.RequestTimeout))
}

func (g *gentlemanImpl) NewDefaultRequest() *gt.Request {
	return g.NewDefaultClient().Request()
}

func (g *gentlemanImpl) NewRequestByConfig(config *GentlemanConfig) *gt.Request {
	return g.NewClientByConfig(config).Request()
}
