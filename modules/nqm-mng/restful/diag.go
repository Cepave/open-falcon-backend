package restful

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/common/diag"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	gin "github.com/gin-gonic/gin"
	json "gopkg.in/bitly/go-simplejson.v0"
)

func health(context *gin.Context) {
	diagRdb := diag.DiagnoseRdb(
		rdb.DbConfig.Dsn,
		rdb.DbFacade.SqlDb,
	)

	jsonResp := json.New()
	jsonResp.Set("rdb", diagRdb)

	jsonHttp := json.New()
	jsonHttp.Set("listening", GinConfig.GetAddress())
	jsonResp.Set("http", jsonHttp)

	context.JSON(http.StatusOK, jsonResp)
}
