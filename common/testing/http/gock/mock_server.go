package gock

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/juju/errors"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gock.v1"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
)

// Functions in namespace for building of *GockConfig
var GockConfigBuilder = &struct {
	NewConfig         func(host string, port uint16) *GockConfig
	NewConfigByRandom func() *GockConfig
}{
	NewConfig: newGockConfig,
	NewConfigByRandom: func() *GockConfig {
		rand.Seed(time.Now().Unix())
		port := rand.Int31n(1000) + 30000
		host := fmt.Sprintf("test-pc%03d.gock.kordan.asshole", rand.Int31n(999)+1)

		return newGockConfig(host, uint16(port))
	},
}

func newGockConfig(host string, port uint16) *GockConfig {
	newConfig := &GockConfig{
		Host: host,
		Port: port,
	}

	url := newConfig.GetUrl()
	newConfig.GentlemanT = &implGentlemanT{url: url}
	newConfig.HttpTest = &implHttptest{mockUrl: url}

	return newConfig
}

// Defines the interface used to ease testing by Gentleman library.
type GentlemanT interface {
	NewClient() *gentleman.Client
	SetupClient(*gentleman.Client) *gentleman.Client
	Plugin() plugin.Plugin
}

// Facade interface used to:
//
// 	1. Mock-up web service with simple configuration
// 	2. Constructs a Gentleman client with configuration of mock
// 	3. Initialize a new *httptest.Server, which could be used to start real server on mock-setup.
type GockConfig struct {
	Host       string
	Port       uint16
	GentlemanT GentlemanT
	HttpTest   tHttp.HttpTest
}

func (c *GockConfig) NewHttpConfig() *client.HttpClientConfig {
	config := client.NewDefaultConfig()
	config.Url = c.GetUrl()
	return config
}

func (c *GockConfig) NewRestfulClientConfig() *oHttp.RestfulClientConfig {
	return &oHttp.RestfulClientConfig{
		HttpClientConfig: c.NewHttpConfig(),
		Plugins: []plugin.Plugin{
			_gentlemanMockPlugin,
		},
	}
}

func (c *GockConfig) GetUrl() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func (c *GockConfig) New() *gock.Request {
	url := c.GetUrl()
	return gock.New(url)
}

// Calls gock.Off()
func (c *GockConfig) Off() {
	logger.Infof("Call gock.Off(): [ %s ]", c.GetUrl())
	gock.Off()
}

// Calls gock.EnableNetworking()
func (c *GockConfig) StartRealNetwork() {
	logger.Infof("Start Gock Real Network[ %s ]", c.GetUrl())
	gock.EnableNetworking()
}

// Calls:
// 	1. gock.Off()
// 	2. gock.DisableNetworking()
func (c *GockConfig) StopRealNetwork() {
	logger.Infof("Stop Gock Real Network[ %s ]", c.GetUrl())
	c.Off()
	gock.DisableNetworking()
}

type implGentlemanT struct {
	url string
}

func (t *implGentlemanT) NewClient() *gentleman.Client {
	return t.SetupClient(gentleman.New())
}
func (t *implGentlemanT) SetupClient(client *gentleman.Client) *gentleman.Client {
	client.BaseURL(t.url).Use(t.Plugin())
	return client
}
func (t *implGentlemanT) Plugin() plugin.Plugin {
	return _gentlemanMockPlugin
}

var _gentlemanMockPlugin = plugin.NewPhasePlugin("before dial", func(ctx *context.Context, h context.Handler) {
	gock.InterceptClient(ctx.Client)
	h.Next(ctx)
})

type implHttptest struct {
	mockUrl string
}

func (self *implHttptest) NewServer(serverConfig *tHttp.FakeServerConfig) *httptest.Server {
	newServer := httptest.NewUnstartedServer(self)
	newServer.Listener = serverConfig.GetListener()
	return newServer
}
func (self *implHttptest) GetHttpHandler() http.Handler {
	return self
}
func (self *implHttptest) ServeHTTP(finalResp http.ResponseWriter, sourceRequest *http.Request) {
	defer func() {
		p := recover()
		if p == nil {
			return
		}

		finalResp.Header().Set("Panic", "Gock handler")
		finalResp.WriteHeader(http.StatusInternalServerError)

		/**
		 * Format error message
		 */
		errorMessage := fmt.Sprintf("%s", p)
		err, ok := p.(error)
		if ok {
			errorMessage = errors.Details(err)
		}
		// :~)

		logger.Errorf("httptest over Gock has error: %s", errorMessage)

		/**
		 * Output error content to JSON
		 */
		jsonBody := ojson.MarshalJSON(
			map[string]interface{}{
				"error": errorMessage,
			},
		)
		finalResp.Write([]byte(jsonBody))
		// :~)
	}()

	client := &http.Client{}

	/**
	 * Re-direct the request to gock client(automatically applying)
	 */
	mockUrl := self.mockUrl + sourceRequest.RequestURI
	finalUrl, err := url.Parse(mockUrl)
	if err != nil {
		err = errors.Annotatef(err, "Cannot parse URL: [%s]", mockUrl)
		panic(errors.Details(err))
	}

	sourceRequest.RequestURI = ""
	sourceRequest.URL = finalUrl
	// :~)

	mockResp, err := client.Do(sourceRequest)
	if err != nil {
		err = errors.Annotatef(err, "Request has error. URL: [%s]", mockUrl)
		panic(errors.Details(err))
	}

	/**
	 * Writes header to real response
	 */
	header := finalResp.Header()
	for k, values := range mockResp.Header {
		for _, value := range values {
			header.Add(k, value)
		}
	}
	finalResp.WriteHeader(mockResp.StatusCode)
	// :~)

	/**
	 * Writes body to real response
	 */
	defer mockResp.Body.Close()
	bytes, err := ioutil.ReadAll(mockResp.Body)
	if err != nil {
		err = errors.Annotate(err, "Load body of mocked response has error")
		panic(errors.Details(err))
	}

	_, err = finalResp.Write(bytes)
	if err != nil {
		err = errors.Annotate(err, "Output response has error")
		panic(errors.Details(err))
	}
	// :~)
}
