package graph

import (
	"github.com/gin-gonic/gin"

	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
)

func InitHttpServices(engine *gin.Engine, mvcBuilder *mvc.MvcBuilder) {
	h := mvcBuilder.BuildHandler

	v1 := engine.Group("/api/v1/graph")

	v1.POST("/endpoint-index/vacuum", h(vacuumEndpointIndex))
}
