package http

import (
	"fmt"
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	"github.com/gin-gonic/gin"
	"gopkg.in/h2non/gentleman.v2"
)

func getHealth(c *gin.Context) {
	req := NewMysqlApiCli().Get().AddPath("/health")
	mysqlInfo := fetchMysqlApiView(req)
	httpInfo := g.Config().Http
	rpcInfo := &g.RpcView{g.Config().Listen}
	heartbeatService := rpc.AgentHeartbeatService

	healthInfo := &g.HealthView{
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

func fetchMysqlApiView(req *gentleman.Request) *g.MysqlApiView {
	var msg string

	res, err := req.Do()
	if err != nil {
		msg = fmt.Sprintf("Err=%v. Body=%s.", err, res.String())
	}
	defer res.Close()

	view := &g.MysqlApiView{
		Address:     req.Context.Request.URL.String(),
		PingResult:  res.StatusCode,
		PingMessage: msg,
	}

	return view
}
