package mvc

import (
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/gin-gonic/gin"
)

var NotFoundOutputBody = OutputBodyFunc(func(c *gin.Context) {
	ogin.JsonNoRouteHandler(c)
})
