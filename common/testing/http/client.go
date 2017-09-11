// Provides both HTTP client and HTTP fake server configuration for testing.
//
// Client Configuration
//
// The "HttpClientConfig" is main configuration object defines tested service of HTTP.
//
// Client Initialization
//
// The "HttpClientConfig" has "NewClient()" and "NewRequest()" to provide out-of-box gentleman client object.
package http

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	json "github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
	checker "gopkg.in/check.v1"
	gt "gopkg.in/h2non/gentleman.v2"

	"github.com/Cepave/open-falcon-backend/common/http/client"
)

// Initialize a client config by flag
//
// 	-http.host(127.0.0.1) - host name of http service
// 	-http.port(80) - port of http service
// 	-http.ssl(false)- whether or not use SSL to test http service
// 	-http.resource("") - the additional resource(sub-path) after host of URL
func NewHttpClientConfigByFlag() *HttpClientConfig {
	var host = flag.String("http.host", "127.0.0.1", "Host of HTTP service to be tested")
	var port = flag.Int("http.port", 80, "Port of HTTP service to be tested")
	var ssl = flag.Bool("http.ssl", false, "Whether or not to use SSL for HTTP service to be tested")
	var resource = flag.String("http.resource", "", "resource for http://<host>:<port/<resource>")

	flag.Parse()

	config := &HttpClientConfig{
		Host:     *host,
		Port:     uint16(*port),
		Ssl:      *ssl,
		Resource: *resource,
	}

	logger.Infof("HTTP URL for testing: %s", config.String())

	return config
}

// The configuration of http client
//
// See "NewHttpClientConfigByFlag()" to initialize this object by flag.
type HttpClientConfig struct {
	Ssl      bool
	Host     string
	Port     uint16
	Resource string
}

// Gets the full URL of tested service
func (self *HttpClientConfig) String() string {
	url := self.hostAndPort()

	if self.Resource != "" {
		url += "/" + self.Resource
	}

	return url
}

func (self *HttpClientConfig) hostAndPort() string {
	schema := "http"
	if self.Ssl {
		schema = "https"
	}

	return fmt.Sprintf("%s://%s:%d", schema, self.Host, self.Port)
}

// Supporting configuration of testing by Gentleman library.
type GentlemanClientConf struct {
	*HttpClientConfig
}

// Consturcts a "*gentleman.Client" object by configuration.
//
// The timeout of request is three seconds.
func (c *GentlemanClientConf) NewClient() *gt.Client {
	gtClient := client.CommonGentleman.NewClientByConfig(
		&client.GentlemanConfig{
			RequestTimeout: time.Duration(3) * time.Second,
		},
	).
		BaseURL(c.String())

	if c.Resource != "" {
		gtClient.Path(c.Resource)
	}

	return gtClient
}

// Consturcts a "*gentleman.Request" object by configuration.
//
// The timeout of request is three seconds.
func (c *GentlemanClientConf) NewRequest() *gt.Request {
	return c.NewClient().Request()
}

// Supporting configuration of testing by Sling(deprecated) library
//
// Deprecated: You should use gentleman library instead.
type SlingClientConf struct {
	*HttpClientConfig
}

func (c *SlingClientConf) NewClient() *sling.Sling {
	client := sling.New().Base(
		c.hostAndPort(),
	)
	if c.Resource != "" {
		client.Path(c.Resource + "/")
	}

	return client
}

// Performs request and reads the body into []byte
func NewResponseResultBySling(slingObj *sling.Sling) *ResponseResult {
	/**
	 * Builds request
	 */
	req, err := slingObj.Request()
	if err != nil {
		panic(err)
	}
	// :~)

	return NewResponseResultByRequest(req)
}

// Performs request and reads the body into []byte
func NewResponseResultByRequest(req *http.Request) *ResponseResult {
	/**
	 * Performs request
	 */
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// :~)

	return NewResponseResultByResponse(resp)
}

func NewResponseResultByResponse(resp *http.Response) *ResponseResult {
	/**
	 * Reads body of response
	 */
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// :~)

	return &ResponseResult{
		Response: resp,
		body:     bodyBytes,
	}
}

type ResponseResult struct {
	Response *http.Response
	body     []byte
}

func (r *ResponseResult) GetBodyAsString() string {
	return string(r.body)
}
func (r *ResponseResult) GetBodyAsJson() *json.Json {
	jsonResult, err := json.NewJson(r.body)
	if err != nil {
		panic(err)
	}

	return jsonResult
}

// Slint with checker
type CheckSlint struct {
	Slint        *sling.Sling
	LastResponse *http.Response

	checker *checker.C
}

// Initialize a checker with slint support
func NewCheckSlint(checker *checker.C, sling *sling.Sling) *CheckSlint {
	return &CheckSlint{
		Slint:   sling,
		checker: checker,
	}
}

// Gets request of slint
func (self *CheckSlint) Request() *http.Request {
	req, err := self.Slint.Request()
	self.checker.Assert(err, checker.IsNil)

	return req
}

// Gets the response for current request
func (self *CheckSlint) GetResponse() *http.Response {
	if self.LastResponse != nil {
		return self.LastResponse
	}

	c := self.checker
	client := &http.Client{}

	var err error
	self.LastResponse, err = client.Do(self.Request())
	c.Assert(err, checker.IsNil)

	return self.LastResponse
}

// Asserts the existing of paging header
func (self *CheckSlint) AssertHasPaging() {
	c := self.checker
	resp := self.GetResponse()

	c.Assert(resp.Header.Get("page-size"), checker.Matches, "\\d+")
	c.Assert(resp.Header.Get("page-pos"), checker.Matches, "\\d+")
	c.Assert(resp.Header.Get("total-count"), checker.Matches, "\\d+")
}

// Gets body as string
//
// The exepcted status is used to get expected status
func (self *CheckSlint) GetStringBody(expectedStatus int) string {
	return string(self.checkAndGetBody(expectedStatus))
}

func (self *CheckSlint) checkAndGetBody(expectedStatus int) []byte {
	c := self.checker

	resp := self.GetResponse()
	defer resp.Body.Close()

	c.Check(resp.StatusCode, checker.Equals, expectedStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if c.Failed() {
		if err != nil {
			c.Fatalf("Read response(ioutil.ReadAll()) has error: %v", err)
		} else {
			c.Fatalf("Status code not match. Response: %s.", bodyBytes)
		}
	}

	return bodyBytes
}

// Gets body as JSON
//
// The exepcted status is used to get expected status
func (self *CheckSlint) GetJsonBody(expectedStatus int) *json.Json {
	c := self.checker

	jsonResult, err := json.NewJson(self.checkAndGetBody(expectedStatus))
	c.Assert(err, checker.IsNil)

	return jsonResult
}
