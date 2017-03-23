package strategy

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
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
