package mockcfg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	mogr := r.Group("/api/v1/nodata")
	mogr.Use(utils.AuthSessionMidd)
	mogr.GET("", GetNoDataList)
	mogr.GET("/:nid", GetNoData)
	mogr.POST("/", CreateNoData)
	mogr.PUT("/", UpdateNoData)
	mogr.DELETE("/:nid", DeleteNoData)
}
