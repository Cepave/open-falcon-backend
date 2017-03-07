package strategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Cepave/open-falcon-backend/modules/api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	strr := r.Group("/api/v1/strategy")
	strr.Use(utils.AuthSessionMidd)
	strr.GET("", GetStrategys)
	strr.GET("/:sid", GetStrategy)
	strr.POST("", CreateStrategy)
	strr.PUT("", UpdateStrategy)
	strr.DELETE("/:sid", DeleteStrategy)
	met := r.Group("/api/v1/metric")
	met.Use(utils.AuthSessionMidd)
	met.GET("tmplist", MetricQuery)
}
