package template

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	tmpr := r.Group("/api/v1/template")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr.GET("", GetTemplates)
	tmpr.POST("", CreateTemplate)
	tmpr.GET("/:tpl_id", GetATemplate)
	tmpr.PUT("", UpdateTemplate)
	tmpr.DELETE("/:tpl_id", DeleteTemplate)
	tmpr.POST("/action", CreateActionToTmplate)
	tmpr.PUT("/action", UpdateActionToTmplate)

	//simple list for ajax use
	tmpr2 := r.Group("/api/v1/template_simple")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr2.GET("", GetTemplatesSimple)
}
