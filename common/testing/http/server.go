package http

import (
	"flag"
	"fmt"

	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"
)

var testGinMode = flag.String("test.gin_mode", gin.DebugMode, "Mode of gin freamework(debug/release/test)")

var testServerHost = flag.String("test.web_host", "0.0.0.0", "Listening Host(0.0.0.0)")
var testServerPort = flag.Uint("test.web_port", 0, "Listening port of web")

var WebTestServer = &struct {
	GetHost func() string
	GetPort func() uint16
	GetUrl func() string
} {
	GetHost: getWebHost,
	GetPort: getWebPort,
	GetUrl: func() string {
		return fmt.Sprintf("http://%s:%d", getWebHost(), getWebPort())
	},
}

var GinTestServer = &struct {
	GoCheckOrSkip            func(c *check.C) bool
	GoCheckStartGinWebServer func(c *check.C, engineFunc GinEngineConfigFunc)
	GoCheckGetConfig         func(c *check.C) *ogin.GinConfig
} {
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
	port := WebTestServer.GetPort()

	if port == 0 {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
	}

	return port > 0
}

func goCheckGetGinConfig(c *check.C) *ogin.GinConfig {
	port := WebTestServer.GetPort()

	if port == 0 {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
		return nil
	}

	return &ogin.GinConfig{
		Mode: *testGinMode,
		Host: WebTestServer.GetHost(),
		Port: port,
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
