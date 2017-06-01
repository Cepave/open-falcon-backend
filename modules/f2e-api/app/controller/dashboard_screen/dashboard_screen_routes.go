package dashboard_screen

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
	authapi.POST("/screen", ScreenCreate)
	authapi.POST("/screen_clone", ScreenClone)
	authapi.PUT("/screen", ScreenUpdate)
	authapi.GET("/screen/:screen_id", ScreenGet)
	authapi.GET("/screens/pid/:pid", ScreenGetsByPid)
	authapi.GET("/screens", ScreenGetsAll)
	authapi.DELETE("/screen/:screen_id", ScreenDelete)
}
