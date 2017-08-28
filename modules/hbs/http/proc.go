package http

import (
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"github.com/gin-gonic/gin"
)

func configProcRoutes(router *gin.Engine) {
	router.GET("/expressions", expressions)
	router.GET("/plugins/:hostname", plugins)
}

func expressions(c *gin.Context) {
	d, err := service.Expressions()
	if err != nil {
		logger.Errorln(err)
	}
	RenderDataJson(c.Writer, d)
}

func plugins(c *gin.Context) {
	d, err := service.Plugins(c.Param("hostname"))
	if err != nil {
		logger.Errorln(err)
	}
	RenderDataJson(c.Writer, d)
}
