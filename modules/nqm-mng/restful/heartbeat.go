package restful

import (
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/service/queue"
	"gopkg.in/gin-gonic/gin.v1"
)

type heartbeatOfAgents []*model.AgentHeartbeat

func (agents *heartbeatOfAgents) Bind(context *gin.Context) {
	ogin.BindJson(context, agents)
}

func agentHeartbeat(
	agents *heartbeatOfAgents,
	q *struct {
		UpdateOnly bool `mvc:"query[update_only]"`
	},
) mvc.OutputBody {
	retBody := rdb.AgentHeartbeat(*agents, q.UpdateOnly)
	return mvc.JsonOutputBody(retBody)
}

func nqmAgentHeartbeat(
	req *model.NqmAgentHeartbeatRequest,
) mvc.OutputBody {
	if rdb.NotNew(req) {
		queue.NqmQueue.Put(req)
	} else {
		rdb.Insert(req)
	}
	r := rdb.SelectByConnId(req.ConnectionId)
	return mvc.JsonOutputBody(r)
}
