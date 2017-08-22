package client

import (
	"time"

	"github.com/h2non/gentleman/plugins/timeout"
	gt "gopkg.in/h2non/gentleman.v2"
)

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
