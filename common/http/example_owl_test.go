package http_test

import (
	"github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
)

// Constructs a client object by set-up configuration.
func ExampleApiService_newClient() {
	httpClientConfig := client.NewDefaultConfig()
	httpClientConfig.Url = "http://some-1.mock.server/"

	restfulConfig := &http.RestfulClientConfig{
		HttpClientConfig: httpClientConfig,
		FromModule:       "query",
	}

	apiService := http.NewApiService(restfulConfig)
	client := apiService.NewClient()

	// The returned request is omitted in the example.
	_ = client.Get()
	// ...
}
