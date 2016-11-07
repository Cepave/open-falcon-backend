package gin

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
)

func OutputJsonIfNotNil(c *gin.Context, body interface{}) {
	if body != nil {
		c.JSON(http.StatusOK, body)
	} else {
		JsonNoRouteHandler(c)
	}
}
