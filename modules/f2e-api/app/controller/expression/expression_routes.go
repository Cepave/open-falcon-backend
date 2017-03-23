package expression

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
	expr := r.Group("/api/v1/expression")
	expr.Use(utils.AuthSessionMidd)
	expr.GET("", GetExpressionList)
	expr.GET("/:eid", GetExpression)
	expr.POST("", CreateExrpession)
	expr.PUT("", UpdateExrpession)
	expr.DELETE("/:eid", DeleteExpression)
}
