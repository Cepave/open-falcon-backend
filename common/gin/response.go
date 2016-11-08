package gin

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"reflect"
)

func OutputJsonIfNotNil(c *gin.Context, checkedObject interface{}) {
	if reflect.ValueOf(checkedObject).IsNil() {
		JsonNoRouteHandler(c)
	} else {
		c.JSON(http.StatusOK, checkedObject)
	}
}
