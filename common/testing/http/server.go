package http

import (
	"flag"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	check "gopkg.in/check.v1"
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
)

var webMode = flag.String("test.gin_mode", gin.DebugMode, "Mode of gin freamework(debug/release/test)")
var webHost = flag.String("test.web_host", "0.0.0.0", "Listening Host(0.0.0.0)")
var webPort = flag.Uint("test.web_port", 0, "Listening port of web")

type GinEngineConfigFunc func(*gin.Engine)

func HasWebConfigOrSkip(c *check.C) bool {
	if *webPort == 0 {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
	}

	return *webPort > 0
}

func GetGinConfig(c *check.C) *ogin.GinConfig {
	if *webPort == 0 {
		c.Skip("Skip mock web testing. Needs \"-test.web_port=<port>\"")
		return nil
	}

	return &ogin.GinConfig {
		Mode: *webMode,
		Host: GetWebHost(),
		Port: uint16(*webPort),
	}
}

func GetWebHost() string {
	if *webHost == "0.0.0.0" {
		return "127.0.0.1"
	}

	return *webHost
}
func GetWebPort() uint16 {
	return uint16(*webPort)
}
func GetWebUrl() string {
	return fmt.Sprintf("http://%s:%d", GetWebHost(), GetWebPort())
}

func StartGinWebServer(c *check.C, engineFunc GinEngineConfigFunc) {
	config := GetGinConfig(c)
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
