/*
This package provides out-of-box configuration for initializing client object to calling of RESTful API services.

Configuration

The "RestfulClientConfig" object defines various properties used to construct calling of OWL service.

Gentleman Client

You could use "ApiService.NewClient()" to get "*gentleman.Client" object.

*/
package http

import (
	gt "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugin"

	"github.com/Cepave/open-falcon-backend/common/http/client"
)

// General configuration to MySqlApi
type RestfulClientConfig struct {
	// Configuration is used to set-up properties of supported library of HTTP client.
	*client.HttpClientConfig
	// If this value is non-empty, this service would add header "From-Module: <FromModule>" in HTTP request.
	FromModule string
	// The "NewClient()" function would USE the "Plugins" to construct "*gentleman.Client" object.
	Plugins []plugin.Plugin
}

// Constructs a new service to API.
func NewApiService(config *RestfulClientConfig) *ApiService {
	return &ApiService{config}
}

// Provides general operation(as service) to RESTful API service
type ApiService struct {
	config *RestfulClientConfig
}

// Constructs a new client object with defined configuration(by RestfulClientConfig)
func (s *ApiService) NewClient() *gt.Client {
	config := s.config

	newClient := client.CommonGentleman.NewClientByConfig(
		&client.GentlemanConfig{
			RequestTimeout: config.RequestTimeout,
		},
	).URL(config.Url)

	if config.FromModule != "" {
		newClient.SetHeader("From-Module", config.FromModule)
	}

	for _, plugin := range s.config.Plugins {
		newClient.Use(plugin)
	}

	return newClient
}
