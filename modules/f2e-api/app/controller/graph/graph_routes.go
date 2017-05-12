package graph

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
	grphapi := r.Group("/api/v1/graph")
	grphapi.Use(utils.AuthSessionMidd)
	grphapi.GET("/endpoint", EndpointRegexpQuery)
	grphapi.GET("/endpoint_counter", EndpointCounterRegexpQuery)
	grphapi.GET("/endpointstr_counter", EndpointStrCounterRegexpQuery)
	grphapi.POST("/history", QueryGraphDrawData)
	owlgraph := r.Group("/api/v1/owlgraph")
	owlgraph.Use(utils.AuthSessionMidd)
	owlgraph.GET("/keyword_search", HostsSearching)
}
