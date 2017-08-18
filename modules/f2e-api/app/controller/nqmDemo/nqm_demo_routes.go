package nqmDemo

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
	nqmd := r.Group("/api/v1/nqm_demo")
	nqmd.GET("/agents", Agents)
	nqmd.GET("/isps", Isps)
	nqmd.GET("/provinces", Provinces)
	nqmd.GET("/targets", Targets)
	nqmd.GET("/pingtasks", PingTasks)
	nqmd.GET("/cities", Cities)
	nqmd.GET("/nametags", NameTags)
	nqmd.GET("/grouptags", GroupTags)
	nqmd2 := r.Group("/api/v1/nqm_demo/email")
	nqmd2.Use(utils.AuthSessionMidd)
	nqmd2.POST("", EmailDemo)
}
