package restful

import (
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"

	"gopkg.in/gin-gonic/gin.v1"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("INFO")
var router *gin.Engine = nil

func InitGin(config *commonGin.GinConfig) {
	if router != nil {
		return
	}

	router = commonGin.NewDefaultJsonEngine(config)

	logger.Infof("Going to start web service. Listen: %s", config)

	initApi()

	go commonGin.StartServiceOrExit(router, config)
}

func initApi() {
	v1 := router.Group("/api/v1")

	v1.GET("/nqm/agents", listAgents)
	v1.GET("/nqm/agent/:agent_id", getAgentById)
	v1.POST("/nqm/agent", addNewAgent)
}
