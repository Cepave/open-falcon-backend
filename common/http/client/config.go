package client

import (
	"time"
)

// Default value of request time out for HTTP client.
const DEFAULT_TIMEOUT = 10 * time.Second

type HttpClientConfig struct {
	Url            string
	RequestTimeout time.Duration
}

// Constructs default configuration by pre-defined values.
func NewDefaultConfig() *HttpClientConfig {
	return &HttpClientConfig{
		RequestTimeout: DEFAULT_TIMEOUT,
	}
}
