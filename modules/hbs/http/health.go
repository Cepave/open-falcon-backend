package http

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"github.com/gin-gonic/gin"
)

func getHealth(c *gin.Context) {
	httpInfo := g.Config().Http
	rpcInfo := &g.RpcView{g.Config().Listen}
	heartbeatService := rpc.AgentHeartbeatService
	mysqlInfo := service.MysqlApiService.GetHealth()

	// Set health check value
	healthcheck := 1
	if mysqlInfo != nil && mysqlInfo.Response != nil {
		healthcheck = mysqlInfo.Response.Rdb.PingResult
	}

	healthInfo := &g.HealthView{
		healthcheck,
		mysqlInfo,
		httpInfo,
		rpcInfo,
		&g.FalconAgentView{
			&g.HeartbeatView{
				CurrentSize:         heartbeatService.CurrentSize(),
				CumulativeReceived:  heartbeatService.CumulativeAgentsPut(),
				CumulativeDropped:   heartbeatService.CumulativeAgentsDropped(),
				CumulativeProcessed: heartbeatService.CumulativeRowsAffected(),
			},
		},
	}
	c.JSON(http.StatusOK, healthInfo)
}
