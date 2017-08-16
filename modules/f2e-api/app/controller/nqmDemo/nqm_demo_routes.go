package nqmDemo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	nqmd := r.Group("/api/v1/nqm_demo")
	nqmd.GET("/agents", Agents)
	nqmd.GET("/isps", Isps)
	nqmd.GET("/provinces", Provinces)
	nqmd.GET("/targets", Targets)
	nqmd.GET("/pingtasks", PingTasks)
	nqmd.GET("/cities", Cities)
	nqmd.GET("/nametags", NameTags)
	nqmd.GET("/grouptags", GroupTags)
}
