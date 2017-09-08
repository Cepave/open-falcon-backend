package http

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	check "gopkg.in/check.v1"

	ogin "github.com/Cepave/open-falcon-backend/common/gin"
)

var testGinMode = flag.String("test.gin_mode", gin.DebugMode, "Mode of gin freamework(debug/release/test)")

var testServerHost = flag.String("test.web_host", "0.0.0.0", "Listening Host(0.0.0.0)")
var testServerPort = flag.Uint("test.web_port", 0, "Listening port of web")

type FakeServerConfig struct {
	Host string
	Port uint16
}

type FakeServer httptest.Server

// Defines the interface used to set-up *net/http/httptest.Server
type HttpTest interface {
	NewServer(serverConfig *FakeServerConfig) *httptest.Server
	GetHttpHandler() http.Handler
}

func (c *FakeServerConfig) GetListener() net.Listener {
	listenerString := fmt.Sprintf("%s:%d", c.Host, c.Port)
	listener, err := net.Listen("tcp", listenerString)

	if err != nil {
		newErr := errors.Annotatef(err, "Cannot create listener[%s]", listenerString)
		panic(errors.Details(newErr))
	}

	return listener
}
func (c *FakeServerConfig) GetUrl() *url.URL {
	urlString := c.GetUrlString()

	urlValue, err := url.Parse(urlString)
	if err != nil {
		newErr := errors.Annotatef(err, "Cannot parse URL string[%s] to \"url.URL\"", urlString)
		panic(errors.Details(newErr))
	}

	return urlValue
}
func (c *FakeServerConfig) GetUrlString() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

// [Deprecated] Try to use testing/http/gock to start a mock server
//
// functions in namespace for getting value of flags for
//
// 	"test.web_host" - The host of web
// 	"test.web_port" - The port of web
//
// The default value of port is "0", hence HasSetting() would be "false"
var WebTestServer = &struct {
	// Gets value of "test.web_host
	GetHost func() string
	// Gets value of "test.web_port
	GetPort func() uint16
	// Gets value of "http://<test.web_host>:<test.web_port>"
	GetUrl     func() string
	HasSetting func() bool
}{
	GetHost: getWebHost,
	GetPort: getWebPort,
	GetUrl: func() string {
		return fmt.Sprintf("http://%s:%d", getWebHost(), getWebPort())
	},
	HasSetting: func() bool {
		return getWebPort() > 0
	},
}

// [Deprecated] Try to use testing/http/gock to start a mock server
//
// functions in namespace for set-up gin server(with GoCheck library)
//
// 	"test.gin_mode" - The mode of Gin framework
// 	"test.web_host" - The host of web
// 	"test.web_port" - The port of web
//
// The default value of port is "0", hence HasSetting() would be "false"
var GinTestServer = &struct {
	GoCheckOrSkip            func(c *check.C) bool
	GoCheckStartGinWebServer func(c *check.C, engineFunc GinEngineConfigFunc)
	GoCheckGetConfig         func(c *check.C) *ogin.GinConfig
}{
	GoCheckOrSkip:            goCheckOrSkip,
	GoCheckStartGinWebServer: goCheckStartGinWebServer,
	GoCheckGetConfig:         goCheckGetGinConfig,
}

type GinEngineConfigFunc func(*gin.Engine)

/* Implementation of WebTestServer */

func getWebHost() string {
	if *testServerHost == "0.0.0.0" {
		return "127.0.0.1"
	}

	return *testServerHost
}

func getWebPort() uint16 {
	return uint16(*testServerPort)
}

/* Implementation of WebTestServer :~) */

/* Implmentaion for GinTestServer */

func goCheckOrSkip(c *check.C) bool {
	result := WebTestServer.HasSetting()

	if !result {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
	}

	return result
}

func goCheckGetGinConfig(c *check.C) *ogin.GinConfig {
	if !WebTestServer.HasSetting() {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
		return nil
	}

	return &ogin.GinConfig{
		Mode: *testGinMode,
		Host: WebTestServer.GetHost(),
		Port: WebTestServer.GetPort(),
	}
}

func goCheckStartGinWebServer(c *check.C, engineFunc GinEngineConfigFunc) {
	config := goCheckGetGinConfig(c)
	if config == nil {
		return
	}

	engine := ogin.NewDefaultJsonEngine(config)

	engineFunc(engine)

	address := config.GetAddress()
	c.Logf("Starting web server at \"%s\"", address)

	go func() {
		err := engine.Run(address)
		c.Assert(err, check.IsNil)
	}()
}

/* Implmentaion for GinTestServer :~) */
