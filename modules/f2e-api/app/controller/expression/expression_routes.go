package expression

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
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
