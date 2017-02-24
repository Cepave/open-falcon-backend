package mvc

import (
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"gopkg.in/gin-gonic/gin.v1"
)

var NotFoundOutputBody = OutputBodyFunc(func(c *gin.Context) {
	ogin.JsonNoRouteHandler(c)
})
