package http

import (
	"fmt"
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
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
	res, err := req.Do()
	view := &g.MysqlApiView{
		Address:    req.Context.Request.URL.String(),
		StatusCode: res.StatusCode,
	}

	if err != nil {
		view.Message += fmt.Sprintf("Err=%v. Body=%s.", err, res.String())
		return view
	}

	h := &apiModel.HealthView{}
	if err = res.JSON(&h); err != nil {
		view.Message += fmt.Sprintf("JSON response Err: %v.", err)
	}
	view.Response = h

	return view
}
