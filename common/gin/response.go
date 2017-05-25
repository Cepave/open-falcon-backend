package gin

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"reflect"
)

// Output JSON if the checkedObject is not nil.
//
// If the checkedObject is nil value, calls "JsonNoRouteHandler"
func OutputJsonIfNotNil(c *gin.Context, checkedObject interface{}) {
	if reflect.ValueOf(checkedObject).IsNil() {
		JsonNoRouteHandler(c)
	} else {
		c.JSON(http.StatusOK, checkedObject)
	}
}
