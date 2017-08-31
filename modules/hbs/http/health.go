package http

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/h2non/gentleman.v2"
)

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Etos(e error) string {
	if e != nil {
		return e.Error()
	}
	return ""
}

func getHealth(c *gin.Context) {
	apiUrl := service.MysqlApiUrl
	heartbeatService := rpc.AgentHeartbeatService
	req := gentleman.New().BaseURL(apiUrl).Path("/health").Head()
	res, err := req.Do()
	rpcInfo := g.Config().Listen
	httpInfo := g.Config().Http
	healthInfo := g.HealthView{
		&g.MysqlApiView{
			Address:     apiUrl,
			PingResult:  Btoi(res.Ok),
			PingMessage: Etos(err),
		},
		httpInfo,
		&g.RpcView{rpcInfo},
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
