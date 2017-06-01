package dashboard_graph

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed
const TMP_GRAPH_FILED_DELIMITER = "|"

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/dashboard")
	authapi.Use(utils.AuthSessionMidd)
	authapi.POST("/tmpgraph", DashboardTmpGraphCreate)
	authapi.POST("/graph", DashboardGraphCreate)
	authapi.POST("/graph_clone", DashboardGraphClone)
	authapi.PUT("/graph", DashboardGraphUpdate)
	authapi.POST("/graph_new_screen", GraphCreateReqDataWithNewScreen)
	authapi.GET("/tmpgraph/:id", DashboardTmpGraphQuery)
	authapi.GET("/graph/:id", DashboardGraphGet)
	authapi.GET("/graphs/screen/:screen_id", DashboardGraphGetsByScreenID)
	authapi.DELETE("/graph/:id", DashboardGraphDelete)
}
