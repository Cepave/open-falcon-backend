package dashboardGraphOwl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masato25/go_email/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	expr := r.Group("/api/v1/dashboard/graph")
	expr.Use(utils.AuthSessionMidd)
	expr.POST("", CreateGraph)
	expr.PUT("", UpdateGraph)
	expr.GET("/:gid", GetGraph)
	expr.DELETE("/:gid", DeleteGraph)
}
