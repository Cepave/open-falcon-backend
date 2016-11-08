package http

import (
	"fmt"
	"flag"
	"net/http"
	"github.com/dghubble/sling"
	"io/ioutil"
	json "github.com/bitly/go-simplejson"
	checker "gopkg.in/check.v1"
)

// Slint with checker
type CheckSlint struct {
	Slint *sling.Sling
	LastResponse *http.Response

	checker *checker.C
}

// Initialize a checker with slint support
func NewCheckSlint(checker *checker.C, sling *sling.Sling) *CheckSlint {
	return &CheckSlint{
		Slint: sling,
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
	c := self.checker

	resp := self.GetResponse()
	defer resp.Body.Close()

	c.Check(resp.StatusCode, checker.Equals, expectedStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if c.Failed() {
		c.Logf("Has error. Response: %s. If ioutil.ReadAll() has error: %v", bodyBytes, err)
		c.FailNow()
	}

	c.Assert(err, checker.IsNil)

	return string(bodyBytes)
}

// Gets body as JSON
//
// The exepcted status is used to get expected status
func (self *CheckSlint) GetJsonBody(expectedStatus int) *json.Json {
	c := self.checker

	resp := self.GetResponse()
	defer resp.Body.Close()

	c.Assert(resp.StatusCode, checker.Equals, expectedStatus)
	jsonResult, err := json.NewFromReader(resp.Body)
	c.Assert(err, checker.IsNil)

	return jsonResult
}

// The configuration of http client
type HttpClientConfig struct {
	Ssl bool
	Host string
	Port uint16
}

// Initialize a client config by flag
//
// 	http_host - host name of http service
// 	http_port - port of http service
// 	http_ssl - whether or not use SSL to test http service
func NewHttpClientConfigByFlag() *HttpClientConfig {
	var host = flag.String("http.host", "127.0.0.1", "Host of HTTP service to be tested")
	var port = flag.Int("http.port", 80, "Port of HTTP service to be tested")
	var ssl = flag.Bool("http.ssl", false, "Whether or not to use SSL for HTTP service to be tested")

	flag.Parse()

	return &HttpClientConfig {
		Host: *host,
		Port: uint16(*port),
		Ssl: *ssl,
	}
}

// Gets the full URL of tested service
func (self *HttpClientConfig) String() string {
	schema := "http"
	if self.Ssl {
		schema = "https"
	}

	return fmt.Sprintf("%s://%s:%d", schema, self.Host, self.Port)
}
