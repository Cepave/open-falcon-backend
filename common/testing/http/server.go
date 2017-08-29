package http

import (
	"flag"
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"
)

var webMode = flag.String("test.gin_mode", gin.DebugMode, "Mode of gin freamework(debug/release/test)")

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

	return &ogin.GinConfig{
		Mode: *webMode,
		Host: GetWebHost(),
		Port: uint16(*webPort),
	}
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
