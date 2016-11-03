package gin

import (
	"fmt"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
)

// This callback function is used to process panic object
type PanicProcessor func(c *gin.Context, panic interface{})

// Builds a gin.HandlerFunc, which is used to handle not-nil object of panic
func BuildJsonPanicProcessor(panicProcessor PanicProcessor) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			panicProcessor(c, p)
		}()

		c.Next()
	}
}

// Type of PanicProcessor, output 500 status with JSON message
func DefaultPanicProcessor(c *gin.Context, panicObject interface{}) {
	c.JSON(
		http.StatusInternalServerError,
		map[string]interface{} {
			"http_status": http.StatusInternalServerError,
			"error_code": -1,
			"error_message": fmt.Sprintf("%v", panicObject),
		},
	)
}

func JsonNoMethodHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{} {
			"http_status": http.StatusMethodNotAllowed,
			"error_code": -1,
			"method": c.Request.Method,
			"uri": c.Request.RequestURI,
		},
	)
}

func JsonNoRouteHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{} {
			"http_status": http.StatusNotFound,
			"error_code": -1,
			"uri": c.Request.RequestURI,
		},
	)
}
