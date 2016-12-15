package gin

import (
	"os"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("INFO")

type GinConfig struct {
	Mode string
	Host string
	Port uint16
}

func (config *GinConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

func (config *GinConfig) String() string {
	return config.GetAddress()
}

// Initialize a router with default JSON response
func NewDefaultJsonEngine(config *GinConfig) *gin.Engine {
	gin.SetMode(config.Mode)

	router := gin.New()
	router.Use(CORSMiddleware)
	router.NoRoute(JsonNoRouteHandler)
	router.NoMethod(JsonNoMethodHandler)
	router.Use(BuildJsonPanicProcessor(DefaultPanicProcessor))

	return router
}

func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func StartServiceOrExit(router *gin.Engine, config *GinConfig) {
	if err := router.Run(config.GetAddress())
		err != nil {
		logger.Errorf("Cannot start web service: %v", err)
		os.Exit(1)
	}
}
