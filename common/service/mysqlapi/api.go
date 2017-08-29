package msyqlapi

import (
	gt "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugin"

	client "github.com/Cepave/open-falcon-backend/common/http/client"
)

// General configuration to MySqlApi
type ApiConfig struct {
	*client.HttpClientConfig
	// If this value is non-empty, this service would add header "from-module: <FromModule>" in HTTP request.
	FromModule string
	Plugins    []plugin.Plugin
}

func NewApiService(config ApiConfig) *ApiService {
	return &ApiService{&config}
}

// Defines general operation for MysqlApiService
type ApiService struct {
	config *ApiConfig
}

func (s *ApiService) NewClient() *gt.Client {
	config := s.config

	newClient := client.CommonGentleman.NewClientByConfig(
		&client.GentlemanConfig{
			RequestTimeout: config.RequestTimeout,
		},
	).BaseURL(config.Url)

	if config.FromModule != "" {
		newClient.AddHeader("from-module", config.FromModule)
	}

	for _, plugin := range s.config.Plugins {
		newClient.Use(plugin)
	}

	return newClient
}
